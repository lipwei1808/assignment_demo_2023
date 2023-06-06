package main

import (
	"context"
	"math/rand"

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
	}

	roomId := req.Message.GetChat()
	err := redisClient.SendMessage(ctx, roomId, message)

	if err != nil {
		return nil, err
	}

	resp.Code, resp.Msg = 0, "success"

	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	roomId := req.GetChat()
	reverse := req.GetReverse()
	cursor := req.GetCursor()
	limit := int64(req.GetLimit())

	end := limit + cursor

	data, err := redisClient.GetMessages(ctx, roomId, cursor, end, reverse)
	if err != nil {
		return nil, err
	}

	resp.Code, resp.Msg = 0, "success"

	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}
