package huron

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/abassian/huron/src/hashgraph"
	"github.com/abassian/huron/src/proxy"
	"github.com/sirupsen/logrus"
)

// SocketHuronProxyServer ...
type SocketHuronProxyServer struct {
	netListener *net.Listener
	rpcServer   *rpc.Server
	handler     proxy.ProxyHandler
	timeout     time.Duration
	logger      *logrus.Logger
}

// NewSocketHuronProxyServer ...
func NewSocketHuronProxyServer(
	bindAddress string,
	handler proxy.ProxyHandler,
	timeout time.Duration,
	logger *logrus.Logger,
) (*SocketHuronProxyServer, error) {

	server := &SocketHuronProxyServer{
		handler: handler,
		timeout: timeout,
		logger:  logger,
	}

	if err := server.register(bindAddress); err != nil {
		return nil, err
	}

	return server, nil
}

func (p *SocketHuronProxyServer) register(bindAddress string) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("State", p)

	p.rpcServer = rpcServer

	l, err := net.Listen("tcp", bindAddress)

	if err != nil {
		return err
	}

	p.netListener = &l

	return nil
}

func (p *SocketHuronProxyServer) listen() error {
	for {
		conn, err := (*p.netListener).Accept()

		if err != nil {
			return err
		}

		go (*p.rpcServer).ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// CommitBlock ...
func (p *SocketHuronProxyServer) CommitBlock(block hashgraph.Block, response *proxy.CommitResponse) (err error) {
	*response, err = p.handler.CommitHandler(block)

	p.logger.WithFields(logrus.Fields{
		"block":    block.Index(),
		"response": response,
		"err":      err,
	}).Debug("HuronProxyServer.CommitBlock")

	return
}

// GetSnapshot ...
func (p *SocketHuronProxyServer) GetSnapshot(blockIndex int, snapshot *[]byte) (err error) {
	*snapshot, err = p.handler.SnapshotHandler(blockIndex)

	if err != nil {
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"block":    blockIndex,
		"snapshot": snapshot,
		"err":      err,
	}).Debug("HuronProxyServer.GetSnapshot")

	return
}

// Restore ...
func (p *SocketHuronProxyServer) Restore(snapshot []byte, stateHash *[]byte) (err error) {
	*stateHash, err = p.handler.RestoreHandler(snapshot)

	if err != nil {
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"state_hash": stateHash,
		"err":        err,
	}).Debug("HuronProxyServer.Restore")

	return
}
