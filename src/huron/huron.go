package huron

import (
	"fmt"
	"os"
	"time"

	"github.com/abassian/huron/src/crypto/keys"
	h "github.com/abassian/huron/src/hashgraph"
	"github.com/abassian/huron/src/net"
	"github.com/abassian/huron/src/node"
	"github.com/abassian/huron/src/peers"
	"github.com/abassian/huron/src/service"
	"github.com/sirupsen/logrus"
)

// Huron is a struct containing the key parts
// of a huron node
type Huron struct {
	Config       *HuronConfig
	Node         *node.Node
	Transport    net.Transport
	Store        h.Store
	Peers        *peers.PeerSet
	GenesisPeers *peers.PeerSet
	Service      *service.Service
}

// NewHuron is a factory method to produce
// a Huron instance.
func NewHuron(config *HuronConfig) *Huron {
	engine := &Huron{
		Config: config,
	}

	return engine
}

// Init initialises the huron engine
func (b *Huron) Init() error {

	if err := b.initPeers(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initPeers")
		return err
	}

	if err := b.initStore(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initStore")
		return err
	}

	if err := b.initTransport(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initTransport")
		return err
	}

	if err := b.initKey(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initKey")
		return err
	}

	if err := b.initNode(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initNode")
		return err
	}

	if err := b.initService(); err != nil {
		b.Config.Logger.WithError(err).Error("huron.go:Init() initService")
		return err
	}

	return nil
}

// Run starts the Huron Node running
func (b *Huron) Run() {
	if b.Service != nil {
		go b.Service.Serve()
	}

	b.Node.Run(true)
}

func (b *Huron) initTransport() error {
	transport, err := net.NewTCPTransport(
		b.Config.BindAddr,
		nil,
		b.Config.MaxPool,
		b.Config.NodeConfig.TCPTimeout,
		b.Config.NodeConfig.JoinTimeout,
		b.Config.Logger,
	)

	if err != nil {
		return err
	}

	b.Transport = transport

	return nil
}

func (b *Huron) initPeers() error {
	if !b.Config.LoadPeers {
		if b.Peers == nil {
			return fmt.Errorf("LoadPeers false, but huron.Peers is nil")
		}

		if b.GenesisPeers == nil {
			return fmt.Errorf("LoadPeers false, but huron.GenesisPeers is nil")
		}

		return nil
	}

	// peers.json
	peerStore := peers.NewJSONPeerSet(b.Config.DataDir, true)

	participants, err := peerStore.PeerSet()
	if err != nil {
		return err
	}

	b.Peers = participants

	// Set Genesis Peer Set from peers.genesis.json
	genesisPeerStore := peers.NewJSONPeerSet(b.Config.DataDir, false)

	genesisParticipants, err := genesisPeerStore.PeerSet()
	if err != nil { // If there is any error, the current peer set is used as the genesis peer set
		b.Config.Logger.Debugf("could not read peers.genesis.json: %v", err)
		b.GenesisPeers = participants
	} else {
		b.GenesisPeers = genesisParticipants
	}

	return nil
}

func (b *Huron) initStore() error {
	if !b.Config.Store {
		b.Config.Logger.Debug("Creating InmemStore")
		b.Store = h.NewInmemStore(b.Config.NodeConfig.CacheSize)
	} else {
		dbPath := b.Config.BadgerDir()

		b.Config.Logger.WithField("path", dbPath).Debug("Creating BadgerStore")

		if !b.Config.NodeConfig.Bootstrap {
			b.Config.Logger.Debug("No Bootstrap")

			backup := backupFileName(dbPath)

			err := os.Rename(dbPath, backup)

			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				b.Config.Logger.Debug("Nothing to backup")
			} else {
				b.Config.Logger.WithField("path", backup).Debug("Created backup")
			}
		}

		b.Config.Logger.WithField("path", dbPath).Debug("Opening BadgerStore")

		dbStore, err := h.NewBadgerStore(b.Config.NodeConfig.CacheSize, dbPath)
		if err != nil {
			return err
		}

		b.Store = dbStore
	}

	return nil
}

func (b *Huron) initKey() error {
	if b.Config.Key == nil {
		simpleKeyfile := keys.NewSimpleKeyfile(b.Config.Keyfile())

		privKey, err := simpleKeyfile.ReadKey()

		if err != nil {
			b.Config.Logger.Warn(fmt.Sprintf("Cannot read private key from file: %v", err))

			privKey, err = keys.GenerateECDSAKey()
			if err != nil {
				b.Config.Logger.Error("Error generating a new ECDSA key")
				return err
			}

			if err := simpleKeyfile.WriteKey(privKey); err != nil {
				b.Config.Logger.Error("Error saving private key", err)
				return err
			}

			b.Config.Logger.Debug("Generated a new private key")
		}

		b.Config.Key = privKey
	}
	return nil
}

func (b *Huron) initNode() error {

	validator := node.NewValidator(b.Config.Key, b.Config.Moniker)

	p, ok := b.Peers.ByID[validator.ID()]
	if ok {
		if p.Moniker != validator.Moniker {
			b.Config.Logger.WithFields(logrus.Fields{
				"json_moniker": p.Moniker,
				"cli_moniker":  validator.Moniker,
			}).Debugf("Using moniker from peers.json file")
			validator.Moniker = p.Moniker
		}
	}

	b.Config.Logger.WithFields(logrus.Fields{
		"genesis_peers": len(b.GenesisPeers.Peers),
		"peers":         len(b.Peers.Peers),
		"id":            validator.ID(),
		"moniker":       validator.Moniker,
	}).Debug("PARTICIPANTS")

	b.Node = node.NewNode(
		&b.Config.NodeConfig,
		validator,
		b.Peers,
		b.GenesisPeers,
		b.Store,
		b.Transport,
		b.Config.Proxy,
	)

	if err := b.Node.Init(); err != nil {
		return fmt.Errorf("failed to initialize node: %s", err)
	}

	return nil
}

func (b *Huron) initService() error {
	if b.Config.ServiceAddr != "" {
		b.Service = service.NewService(b.Config.ServiceAddr, b.Node, b.Config.Logger)
	}
	return nil
}

// backupFileName implements the naming convention for database backups:
// badger_db--UTC--<created_at UTC ISO8601>
func backupFileName(base string) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("%s--UTC--%s", base, toISO8601(ts))
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}
