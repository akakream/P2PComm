package datastore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akakream/sailorsailor/p2p"
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

func NewDatastore(ctx context.Context, pubsubClient *p2p.LibP2PClient, datapath string) (*Datastore, error) {
	store, err := badger.NewDatastore(datapath, &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}

	topicName := "globaldb-topic"
	crdt, err := setupCRDT(ctx, pubsubClient, topicName, store)
	if err != nil {
		log.Fatal(err)
	}

	return &Datastore{
		datapath:  datapath,
		Store:     store,
		Crdt:      crdt,
		topicName: topicName,
	}, nil
}

func (d *Datastore) Shutdown() {
	d.Store.Close()
	d.Crdt.Close()
}

func setupCRDT(ctx context.Context, pubsubClient *p2p.LibP2PClient, topicName string, store *badger.Datastore) (*crdt.Datastore, error) {
	pubsubBC, err := crdt.NewPubSubBroadcaster(ctx, pubsubClient.Ps, topicName)
	if err != nil {
		return nil, err
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
	crdt, err := crdt.New(store, ds.NewKey("crdt"), nil, pubsubBC, opts)
	if err != nil {
		return nil, err
	}

	return crdt, nil
}
