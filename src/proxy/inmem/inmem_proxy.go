package inmem

import (
	hg "github.com/abassian/huron/src/hashgraph"
	"github.com/abassian/huron/src/proxy"
	"github.com/sirupsen/logrus"
)

//InmemProxy implements the AppProxy interface natively
type InmemProxy struct {
	handler  proxy.ProxyHandler
	submitCh chan []byte
	logger   *logrus.Logger
}

// NewInmemProxy instantiates an InmemProxy from a set of handlers.
// If no logger, a new one is created
func NewInmemProxy(handler proxy.ProxyHandler,
	logger *logrus.Logger) *InmemProxy {

	if logger == nil {
		logger = logrus.New()

		logger.Level = logrus.DebugLevel
	}

	return &InmemProxy{
		handler:  handler,
		submitCh: make(chan []byte),
		logger:   logger,
	}
}

/*******************************************************************************
* SubmitTx                                                                     *
*******************************************************************************/

//SubmitTx is called by the App to submit a transaction to Huron
func (p *InmemProxy) SubmitTx(tx []byte) {
	//have to make a copy, or the tx will be garbage collected and weird stuff
	//happens in transaction pool
	t := make([]byte, len(tx), len(tx))

	copy(t, tx)

	p.submitCh <- t
}

/*******************************************************************************
* Implement AppProxy Interface                                                 *
*******************************************************************************/

//SubmitCh returns the channel of raw transactions
func (p *InmemProxy) SubmitCh() chan []byte {
	return p.submitCh
}

//CommitBlock calls the commitHandler
func (p *InmemProxy) CommitBlock(block hg.Block) (proxy.CommitResponse, error) {
	commitResponse, err := p.handler.CommitHandler(block)

	p.logger.WithFields(logrus.Fields{
		"index":          block.Index(),
		"round_received": block.RoundReceived(),
		"frame_hash":     block.FrameHash(),
		"peers_hash":     block.PeersHash(),
		"state_hash":     block.StateHash(),
		"txs":            len(block.Transactions()),
		"response":       commitResponse,
		"err":            err,
	}).Debug("InmemProxy.CommitBlock")

	return commitResponse, err
}

//GetSnapshot calls the snapshotHandler
func (p *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	snapshot, err := p.handler.SnapshotHandler(blockIndex)

	p.logger.WithFields(logrus.Fields{
		"block":    blockIndex,
		"snapshot": snapshot,
		"err":      err,
	}).Debug("InmemProxy.GetSnapshot")

	return snapshot, err
}

//Restore calls the restoreHandler
func (p *InmemProxy) Restore(snapshot []byte) error {
	stateHash, err := p.handler.RestoreHandler(snapshot)

	p.logger.WithFields(logrus.Fields{
		"state_hash": stateHash,
		"err":        err,
	}).Debug("InmemProxy.Restore")

	return err
}
