package huron

import (
	"fmt"
	"time"

	"github.com/abassian/huron/src/proxy"
	"github.com/sirupsen/logrus"
)

// SocketHuronProxy ...
type SocketHuronProxy struct {
	nodeAddress string
	bindAddress string

	handler proxy.ProxyHandler

	client *SocketHuronProxyClient
	server *SocketHuronProxyServer
}

// NewSocketHuronProxy ...
func NewSocketHuronProxy(
	nodeAddr string,
	bindAddr string,
	handler proxy.ProxyHandler,
	timeout time.Duration,
	logger *logrus.Logger,
) (*SocketHuronProxy, error) {

	if logger == nil {
		logger = logrus.New()

		logger.Level = logrus.DebugLevel
	}

	client := NewSocketHuronProxyClient(nodeAddr, timeout)

	server, err := NewSocketHuronProxyServer(bindAddr, handler, timeout, logger)

	if err != nil {
		return nil, err
	}

	proxy := &SocketHuronProxy{
		nodeAddress: nodeAddr,
		bindAddress: bindAddr,
		handler:     handler,
		client:      client,
		server:      server,
	}

	go proxy.server.listen()

	return proxy, nil
}

// SubmitTx ...
func (p *SocketHuronProxy) SubmitTx(tx []byte) error {
	ack, err := p.client.SubmitTx(tx)

	if err != nil {
		return err
	}

	if !*ack {
		return fmt.Errorf("Failed to deliver transaction to Huron")
	}

	return nil
}
