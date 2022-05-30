package agent

import (
	"context"
	"io"
	"log"

	"github.com/recode-sh/agent/proto"
)

type InitInstanceStream interface {
	Recv() (*proto.InitInstanceReply, error)
}

func (c Client) InitInstance(
	initInstanceRequest *proto.InitInstanceRequest,
	streamHandler func(stream InitInstanceStream) error,
) error {

	return c.Execute(func(agentGRPCClient proto.AgentClient) error {
		initInstanceStream, err := agentGRPCClient.InitInstance(
			context.TODO(),
			initInstanceRequest,
		)

		if err != nil {
			return err
		}

		return streamHandler(initInstanceStream)
	})
}

func InitInstanceDefaultStreamHandler(stream InitInstanceStream) error {
	for {
		initInstanceReply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Println(initInstanceReply.LogLine)
	}

	return nil
}
