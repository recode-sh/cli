package agent

import (
	"github.com/recode-sh/agent/proto"
	"github.com/recode-sh/cli/internal/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientInterface interface {
	InitInstance(
		initInstanceRequest *proto.InitInstanceRequest,
		streamHandler func(stream InitInstanceStream) error,
	) error

	BuildAndStartDevEnv(
		startDevEnvRequest *proto.BuildAndStartDevEnvRequest,
		streamHandler func(stream BuildAndStartDevEnvStream) error,
	) error
}

type ClientConfig struct {
	ServerRootUser           string
	ServerSSHPrivateKeyBytes []byte
	ServerAddr               string
	LocalAddr                string
	RemoteAddrProtocol       string
	RemoteAddr               string
}

type Client struct {
	config ClientConfig
}

func NewClient(config ClientConfig) Client {
	return Client{
		config: config,
	}
}

func (c Client) Execute(fnToExec func(agentGRPCClient proto.AgentClient) error) error {
	portForwarderReadyChan := make(chan ssh.PortForwarderReadyResp)
	portForwarderRespChan := make(chan error)
	portForwarder := ssh.NewPortForwarder()

	// Open an SSH tunnel to the GRPC server
	// from "localAddr" to "remoteAddr" inside server ("serverAddr")
	go func() {
		portForwarderRespChan <- portForwarder.Forward(
			portForwarderReadyChan,
			c.config.ServerSSHPrivateKeyBytes,
			c.config.ServerRootUser,
			c.config.ServerAddr,
			c.config.LocalAddr,
			c.config.RemoteAddrProtocol,
			c.config.RemoteAddr,
		)
	}()

	portForwarderReadyResp := <-portForwarderReadyChan

	if portForwarderReadyResp.Error != nil {
		return portForwarderReadyResp.Error
	}

	grpcRespChan := make(chan error)

	go func() {
		grpcConn, err := grpc.Dial(
			portForwarderReadyResp.LocalAddr,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)

		if err != nil {
			grpcRespChan <- err
			return
		}

		defer grpcConn.Close()

		agentGRPCClient := proto.NewAgentClient(grpcConn)

		grpcRespChan <- fnToExec(agentGRPCClient)
	}()

	select {
	case portForwarderErr := <-portForwarderRespChan:
		return portForwarderErr
	case grpcErr := <-grpcRespChan:
		return grpcErr
	}
}
