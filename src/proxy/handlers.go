package proxy

import "github.com/abassian/huron/src/hashgraph"

/*
These types are exported and need to be implemented and used by the calling
application.
*/

// ProxyHandler ...
type ProxyHandler interface {
	//CommitHandler is called when Huron commits a block to the application. It
	//returns the state hash resulting from applying the block's transactions to the
	//state
	CommitHandler(block hashgraph.Block) (response CommitResponse, err error)

	//SnapshotHandler is called by Huron to retrieve a snapshot corresponding to a
	//particular block
	SnapshotHandler(blockIndex int) (snapshot []byte, err error)

	//RestoreHandler is called by Huron to restore the application to a specific
	//state
	RestoreHandler(snapshot []byte) (stateHash []byte, err error)
}
