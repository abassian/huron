package dummy

import (
	"time"

	socket "github.com/abassian/huron/src/proxy/socket/huron"
	"github.com/sirupsen/logrus"
)

//DummySocketClient is a socket implementation of the dummy app. Huron and the
//app run in separate processes and communicate through TCP sockets using
//a SocketHuronProxy and a SocketAppProxy.
type DummySocketClient struct {
	state       *State
	huronProxy *socket.SocketHuronProxy
	logger      *logrus.Logger
}

//NewDummySocketClient instantiates a DummySocketClient and starts the
//SocketHuronProxy
func NewDummySocketClient(clientAddr string, nodeAddr string, logger *logrus.Logger) (*DummySocketClient, error) {
	state := NewState(logger)

	huronProxy, err := socket.NewSocketHuronProxy(nodeAddr, clientAddr, state, 1*time.Second, logger)

	if err != nil {
		return nil, err
	}

	client := &DummySocketClient{
		state:       state,
		huronProxy: huronProxy,
		logger:      logger,
	}

	return client, nil
}

//SubmitTx sends a transaction to Huron via the SocketProxy
func (c *DummySocketClient) SubmitTx(tx []byte) error {
	return c.huronProxy.SubmitTx(tx)
}
