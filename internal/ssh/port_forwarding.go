package ssh

import (
	"io"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type PortForwarder struct{}

func NewPortForwarder() PortForwarder {
	return PortForwarder{}
}

type PortForwarderReadyResp struct {
	Error     error
	LocalAddr string
}

func (p PortForwarder) Forward(
	onReadyChan chan<- PortForwarderReadyResp,
	privateKeyBytes []byte,
	user string,
	serverAddr string,
	localAddr string,
	remoteAddrProtocol string,
	remoteAddr string,
) error {

	sshConfig, err := p.buildSSHConfig(
		8*time.Second,
		user,
		privateKeyBytes,
	)

	if err != nil {
		onReadyChan <- PortForwarderReadyResp{
			Error: err,
		}
		return nil
	}

	// Establish connection with server through SSH
	serverSSHConn, err := ssh.Dial("tcp", serverAddr, sshConfig)

	if err != nil {
		onReadyChan <- PortForwarderReadyResp{
			Error: err,
		}
		return nil
	}

	defer serverSSHConn.Close()

	// Establish connection with remoteAddr from server
	remoteConn, err := serverSSHConn.Dial(remoteAddrProtocol, remoteAddr)

	if err != nil {
		onReadyChan <- PortForwarderReadyResp{
			Error: err,
		}
		return nil
	}

	// Start local TCP server to forward traffic to remote connection
	localTCPServer, err := net.Listen("tcp", localAddr)

	if err != nil {
		onReadyChan <- PortForwarderReadyResp{
			Error: err,
		}
		return nil
	}

	defer localTCPServer.Close()

	onReadyChan <- PortForwarderReadyResp{
		LocalAddr: localTCPServer.Addr().String(),
	}

	localConn, err := localTCPServer.Accept()

	if err != nil {
		return err
	}

	return p.forwardLocalToRemoteConn(
		localConn,
		remoteConn,
	)
}

// Get ssh client config for our connection
func (p PortForwarder) buildSSHConfig(
	connTimeout time.Duration,
	user string,
	privateKeyBytes []byte,
) (*ssh.ClientConfig, error) {

	parsedPrivateKey, err := ssh.ParsePrivateKey(privateKeyBytes)

	if err != nil {
		return nil, err
	}

	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(parsedPrivateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         connTimeout,
	}

	return &config, nil
}

// Handle local TCP server connections and tunnel data to the remote server
func (p PortForwarder) forwardLocalToRemoteConn(
	localConn net.Conn,
	remoteConn net.Conn,
) error {

	defer func() {
		localConn.Close()
		remoteConn.Close()
	}()

	remoteConnRespChan := make(chan error, 1)
	localConnRespChan := make(chan error, 1)

	// Forward remote -> local
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		remoteConnRespChan <- err
	}()

	// Forward local -> remote
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		localConnRespChan <- err
	}()

	select {
	case remoteConnErr := <-remoteConnRespChan:
		return remoteConnErr
	case localConnErr := <-localConnRespChan:
		return localConnErr
	}
}
