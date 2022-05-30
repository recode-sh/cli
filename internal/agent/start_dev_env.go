package agent

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/recode-sh/agent/proto"
	"github.com/recode-sh/cli/internal/constants"
)

type BuildAndStartDevEnvStream interface {
	Recv() (*proto.BuildAndStartDevEnvReply, error)
}

func (c Client) BuildAndStartDevEnv(
	startDevEnvRequest *proto.BuildAndStartDevEnvRequest,
	streamHandler func(stream BuildAndStartDevEnvStream) error,
) error {

	return c.Execute(func(agentGRPCClient proto.AgentClient) error {
		startDevEnvStream, err := agentGRPCClient.BuildAndStartDevEnv(
			context.TODO(),
			startDevEnvRequest,
		)

		if err != nil {
			return err
		}

		return streamHandler(startDevEnvStream)
	})
}

func BuildAndStartDevEnvDefaultStreamHandler(stream BuildAndStartDevEnvStream) error {
	for {
		startDevEnvReply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(startDevEnvReply.LogLineHeader) > 0 {
			bold := constants.Bold
			fmt.Println(bold("[" + startDevEnvReply.LogLineHeader + "]\n"))
		}

		if len(startDevEnvReply.LogLine) > 0 {
			log.Println(startDevEnvReply.LogLine)
		}
	}

	return nil
}
