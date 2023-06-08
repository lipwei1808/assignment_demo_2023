package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()

	message := &Message{
		Sender:    req.Message.GetSender(),
		Message:   req.Message.GetText(),
		Timestamp: req.Message.GetSendTime(),
		Header:    req.Message.GetHeader(),
	}

	roomId, e := getRoomId(req.Message.GetChat())
	fmt.Printf("Inside RPC Server, sending message in %q\n", roomId)
	if e != nil {
		fmt.Println("Error in sending messages")
		return nil, e
	}

	err := redisClient.SendMessage(ctx, roomId, message)

	if err != nil {
		return nil, err
	}

	resp.Code, resp.Msg = 0, "success"

	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	roomId, err := getRoomId(req.GetChat())
	if err != nil {
		return nil, err
	}
	reverse := req.GetReverse()
	cursor := req.GetCursor()
	limit := int64(req.GetLimit())

	end := limit + cursor

	data, err := redisClient.GetMessages(ctx, roomId, cursor, end, reverse)
	if err != nil {
		return nil, err
	}

	var counter int64 = 0
	var hasMore bool
	var nextCursor int64
	var messages []*rpc.Message
	for _, row := range data {
		if counter == limit {
			hasMore = true
			nextCursor = end
			break
		}
		fmt.Println(row.Message, row.Sender, row.Timestamp, row.Header)
		temp := &rpc.Message{
			Chat:     roomId,
			Text:     row.Message,
			Sender:   row.Sender,
			SendTime: row.Timestamp,
			Header:   row.Header,
		}
		messages = append(messages, temp)
		counter += 1
	}

	resp.Code, resp.Msg = 0, "success"
	resp.HasMore = &hasMore
	resp.Messages = messages
	resp.NextCursor = &nextCursor
	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}

func getRoomId(chat string) (string, error) {
	persons := strings.Split(chat, ":")
	if len(persons) != 2 {
		return "", errors.New("invalid id")
	}
	sort.Strings(persons)
	return persons[0] + ":" + persons[1], nil
}
