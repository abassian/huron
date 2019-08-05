package socket

import (
	"reflect"
	"testing"
	"time"

	"github.com/abassian/huron/src/common"
	"github.com/abassian/huron/src/hashgraph"
	"github.com/abassian/huron/src/peers"
	"github.com/abassian/huron/src/proxy"
	aproxy "github.com/abassian/huron/src/proxy/socket/app"
	bproxy "github.com/abassian/huron/src/proxy/socket/huron"
	"github.com/sirupsen/logrus"
)

type TestHandler struct {
	blocks     []hashgraph.Block
	blockIndex int
	snapshot   []byte
	logger     *logrus.Logger
}

func (p *TestHandler) CommitHandler(block hashgraph.Block) (proxy.CommitResponse, error) {
	p.logger.Debug("CommitBlock")

	p.blocks = append(p.blocks, block)

	receipts := []hashgraph.InternalTransactionReceipt{}
	for _, it := range block.InternalTransactions() {
		receipts = append(receipts, it.AsAccepted())
	}

	response := proxy.CommitResponse{
		StateHash:                   []byte("statehash"),
		InternalTransactionReceipts: receipts,
	}

	return response, nil
}

func (p *TestHandler) SnapshotHandler(blockIndex int) ([]byte, error) {
	p.logger.Debug("GetSnapshot")

	p.blockIndex = blockIndex

	return []byte("snapshot"), nil
}

func (p *TestHandler) RestoreHandler(snapshot []byte) ([]byte, error) {
	p.logger.Debug("RestoreSnapshot")

	p.snapshot = snapshot

	return []byte("statehash"), nil
}

func NewTestHandler(t *testing.T) *TestHandler {
	logger := common.NewTestLogger(t)

	return &TestHandler{
		blocks:     []hashgraph.Block{},
		blockIndex: 0,
		snapshot:   []byte{},
		logger:     logger,
	}
}

func TestSocketProxyServer(t *testing.T) {
	clientAddr := "127.0.0.1:6990"
	proxyAddr := "127.0.0.1:6991"

	appProxy, err := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))

	if err != nil {
		t.Fatalf("Cannot create SocketAppProxy: %s", err)
	}

	submitCh := appProxy.SubmitCh()

	tx := []byte("the test transaction")

	// Listen for a request
	go func() {
		select {
		case st := <-submitCh:
			// Verify the command
			if !reflect.DeepEqual(st, tx) {
				t.Fatalf("tx mismatch: %#v %#v", tx, st)
			}

		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// now client part connecting to RPC service
	// and calling methods
	huronProxy, err := bproxy.NewSocketHuronProxy(proxyAddr, clientAddr, NewTestHandler(t), 1*time.Second, common.NewTestLogger(t))

	if err != nil {
		t.Fatal(err)
	}

	err = huronProxy.SubmitTx(tx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSocketProxyClient(t *testing.T) {
	clientAddr := "127.0.0.1:6992"
	proxyAddr := "127.0.0.1:6993"

	logger := common.NewTestLogger(t)

	//create app proxy
	appProxy, err := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, logger)
	if err != nil {
		t.Fatalf("Cannot create SocketAppProxy: %s", err)
	}

	handler := NewTestHandler(t)

	//create huron proxy
	_, err = bproxy.NewSocketHuronProxy(proxyAddr, clientAddr, handler, 1*time.Second, logger)
	if err != nil {
		t.Fatalf("Cannot create SocketHuronProxy: %s", err)
	}

	transactions := [][]byte{
		[]byte("tx 1"),
		[]byte("tx 2"),
		[]byte("tx 3"),
	}

	internalTransactions := []hashgraph.InternalTransaction{
		hashgraph.NewInternalTransaction(hashgraph.PEER_ADD, *peers.NewPeer("node0", "paris", "")),
		hashgraph.NewInternalTransaction(hashgraph.PEER_REMOVE, *peers.NewPeer("node1", "london", "")),
	}

	block := hashgraph.NewBlock(0, 1, []byte{}, []*peers.Peer{}, transactions, internalTransactions)

	expectedStateHash := []byte("statehash")
	expectedSnapshot := []byte("snapshot")

	commitResponse, err := appProxy.CommitBlock(*block)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(block.Body, handler.blocks[0].Body) {
		t.Fatalf("block should be \n%#v\n, not \n%#v\n", *block, handler.blocks[0])
	}

	if !reflect.DeepEqual(commitResponse.StateHash, expectedStateHash) {
		t.Fatalf("StateHash should be %v, not %v", expectedStateHash, commitResponse.StateHash)
	}

	snapshot, err := appProxy.GetSnapshot(block.Index())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(block.Index(), handler.blockIndex) {
		t.Fatalf("blockIndex should be %v, not %v", block.Index(), handler.blockIndex)
	}

	if !reflect.DeepEqual(snapshot, expectedSnapshot) {
		t.Fatalf("Snapshot should be %v, not %v", expectedSnapshot, snapshot)
	}

	err = appProxy.Restore(snapshot)
	if err != nil {
		t.Fatalf("Error restoring snapshot: %v", err)
	}

	if !reflect.DeepEqual(expectedSnapshot, handler.snapshot) {
		t.Fatalf("snapshot should be %v, not %v", expectedSnapshot, handler.snapshot)
	}
}
