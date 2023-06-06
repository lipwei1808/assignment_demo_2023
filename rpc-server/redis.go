package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (c *RedisClient) InitRedisClient(ctx context.Context) error {
	fmt.Println("Starting Redis Client")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	c.cli = client
	return nil
}

func (c *RedisClient) SendMessage(ctx context.Context, roomId string, msg *Message) error {
	data, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	member := &redis.Z{
		Score:  float64(msg.Timestamp),
		Member: data,
	}

	_, err = c.cli.ZAdd(ctx, roomId, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessages(ctx context.Context, roomId string, start, stop int64, reverse bool) ([]*Message, error) {
	var messages []*Message
	var data []string
	var err error

	if !reverse {
		data, err = c.cli.ZRange(ctx, roomId, start, stop).Result()
	} else {
		data, err = c.cli.ZRevRange(ctx, roomId, start, stop).Result()
	}

	if err != nil {
		return nil, err
	}

	for _, row := range data {
		temp := &Message{}

		err := json.Unmarshal([]byte(row), temp)
		if err != nil {
			return nil, err
		}

		messages = append(messages, temp)
	}

	return messages, nil
}
