package huron

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	bkeys "github.com/abassian/huron/src/crypto/keys"
	"github.com/abassian/huron/src/peers"
)

func TestInitStore(t *testing.T) {
	os.RemoveAll("test_data")
	os.Mkdir("test_data", os.ModeDir|0777)
	defer os.RemoveAll("test_data")

	conf := NewDefaultConfig()
	conf.DataDir = "test_data"
	conf.Store = true
	conf.NodeConfig.Bootstrap = false

	jsonPeerSet := peers.NewJSONPeerSet("test_data", true)

	keys := map[string]*ecdsa.PrivateKey{}
	peerSlice := []*peers.Peer{}
	for i := 0; i < 3; i++ {
		key, _ := bkeys.GenerateECDSAKey()
		peer := &peers.Peer{
			NetAddr:   fmt.Sprintf("addr%d", i),
			PubKeyHex: bkeys.PublicKeyHex(&key.PublicKey),
			Moniker:   fmt.Sprintf("peer%d", i),
		}
		peerSlice = append(peerSlice, peer)
		keys[peer.NetAddr] = key
	}

	newPeerSet := peers.NewPeerSet(peerSlice)
	newPeerSlice := newPeerSet.Peers

	if err := jsonPeerSet.Write(newPeerSlice); err != nil {
		t.Fatalf("err: %v", err)
	}

	huron := NewHuron(conf)

	if err := huron.initStore(); err != nil {
		t.Fatal(err)
	}

	huron2 := NewHuron(conf)

	if err := huron2.initStore(); err != nil {
		t.Fatal(err)
	}

	// check that huron2 created a backup
	files, err := ioutil.ReadDir("test_data")
	if err != nil {
		t.Fatal(err)
	}
	dbFiles := []string{}
	for _, f := range files {
		if strings.Contains(f.Name(), "badger_db") {
			dbFiles = append(dbFiles, f.Name())
		}
	}
	if len(dbFiles) != 2 {
		t.Fatalf("initStore should have created a new db file")
	}
}
