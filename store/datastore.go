package store

import (
	"context"
	"fmt"
	"time"

	"github.com/akakream/sailorsailor/p2p"
	ipfslite "github.com/hsanjuan/ipfs-lite"
	ds "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger"
	crdt "github.com/ipfs/go-ds-crdt"
)

type Datastore struct {
	Store     *badger.Datastore
	Crdt      *crdt.Datastore
	datapath  string
	topicName string
}

func NewDatastore(ctx context.Context, datapath string) (*Datastore, error) {
	store, err := badger.NewDatastore(datapath, &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return &Datastore{
		datapath:  datapath,
		Store:     store,
		topicName: "globaldb-topic",
	}, nil
}

func (d *Datastore) Shutdown() {
	d.Store.Close()
	d.Crdt.Close()
}

func (d *Datastore) SetupCRDT(ctx context.Context, pubsubClient *p2p.LibP2PClient, liteipfs *ipfslite.Peer) error {
	pubsubBC, err := crdt.NewPubSubBroadcaster(ctx, pubsubClient.Ps, d.topicName)
	if err != nil {
		return err
	}

	opts := crdt.DefaultOptions()
	opts.RebroadcastInterval = 5 * time.Second
	opts.PutHook = func(k ds.Key, v []byte) {
		fmt.Printf("Added: [%s] -> %s\n", k, string(v))

	}
	opts.DeleteHook = func(k ds.Key) {
		fmt.Printf("Removed: [%s]\n", k)
	}

	// Add ipfs-lite
	crdt, err := crdt.New(d.Store, ds.NewKey("crdt"), liteipfs, pubsubBC, opts)
	if err != nil {
		return err
	}
	d.Crdt = crdt

	return nil
}
