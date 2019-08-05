package huron

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

// SocketHuronProxyClient ...
type SocketHuronProxyClient struct {
	nodeAddr string
	timeout  time.Duration
	rpc      *rpc.Client
}

// NewSocketHuronProxyClient ...
func NewSocketHuronProxyClient(nodeAddr string, timeout time.Duration) *SocketHuronProxyClient {
	return &SocketHuronProxyClient{
		nodeAddr: nodeAddr,
		timeout:  timeout,
	}
}

func (p *SocketHuronProxyClient) getConnection() error {
	if p.rpc == nil {
		conn, err := net.DialTimeout("tcp", p.nodeAddr, p.timeout)

		if err != nil {
			return err
		}

		p.rpc = jsonrpc.NewClient(conn)
	}

	return nil
}

// SubmitTx ...
func (p *SocketHuronProxyClient) SubmitTx(tx []byte) (*bool, error) {
	if err := p.getConnection(); err != nil {
		return nil, err
	}

	var ack bool

	err := p.rpc.Call("Huron.SubmitTx", tx, &ack)

	if err != nil {
		p.rpc = nil

		return nil, err
	}

	return &ack, nil
}
