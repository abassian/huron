package mobile

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abassian/huron/src/huron"
	"github.com/abassian/huron/src/crypto/keys"
	"github.com/abassian/huron/src/node"
	"github.com/abassian/huron/src/peers"
	"github.com/abassian/huron/src/proxy"
	"github.com/abassian/huron/src/proxy/inmem"
	"github.com/sirupsen/logrus"
)

// Node ...
type Node struct {
	nodeID uint32
	node   *node.Node
	proxy  proxy.AppProxy
	logger *logrus.Logger
}

// New initializes Node struct
func New(privKey string,
	nodeAddr string,
	jsonPeers string,
	commitHandler CommitHandler,
	exceptionHandler ExceptionHandler,
	config *MobileConfig) *Node {

	huronConfig := config.toHuronConfig()

	huronConfig.Logger.WithFields(logrus.Fields{
		"nodeAddr": nodeAddr,
		"peers":    jsonPeers,
		"config":   fmt.Sprintf("%v", config),
	}).Debug("New Mobile Node")

	huronConfig.BindAddr = nodeAddr

	//Check private key
	keyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Failed to decode private key bytes: %s", err))
		return nil
	}

	key, err := keys.ParsePrivateKey(keyBytes)
	if err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Failed to parse private key: %s", err))
		return nil
	}

	huronConfig.Key = key

	// Decode the peers
	var ps []*peers.Peer
	dec := json.NewDecoder(strings.NewReader(jsonPeers))
	if err := dec.Decode(&ps); err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Failed to parse PeerSet: %s", err))
		return nil
	}

	peerSet := peers.NewPeerSet(ps)

	huronConfig.LoadPeers = false

	//mobileApp implements the ProxyHandler interface, and we use it to
	//instantiate an InmemProxy
	mobileApp := newMobileApp(commitHandler, exceptionHandler, huronConfig.Logger)
	huronConfig.Proxy = inmem.NewInmemProxy(mobileApp, huronConfig.Logger)

	engine := huron.NewHuron(huronConfig)

	engine.Peers = peerSet
	engine.GenesisPeers = peerSet

	if err := engine.Init(); err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Cannot initialize engine: %s", err))
		return nil
	}

	return &Node{
		node:   engine.Node,
		proxy:  huronConfig.Proxy,
		nodeID: engine.Node.GetID(),
		logger: huronConfig.Logger,
	}
}

// Run ...
func (n *Node) Run(async bool) {
	if async {
		n.node.RunAsync(true)
	} else {
		n.node.Run(true)
	}
}

// Leave ...
func (n *Node) Leave() {
	n.node.Leave()
}

// Shutdown ...
func (n *Node) Shutdown() {
	n.node.Shutdown()
}

// SubmitTx ...
func (n *Node) SubmitTx(tx []byte) {
	//have to make a copy or the tx will be garbage collected and weird stuff
	//happens in transaction pool
	t := make([]byte, len(tx), len(tx))
	copy(t, tx)
	n.proxy.SubmitCh() <- t
}

// GetPeers ...
func (n *Node) GetPeers() string {
	peers := n.node.GetPeers()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(peers); err != nil {
		return ""
	}

	return buf.String()
}

// GetStats ...
func (n *Node) GetStats() string {
	stats := n.node.GetStats()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(stats); err != nil {
		return ""
	}

	return buf.String()
}
