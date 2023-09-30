package store

import (
	"context"
	"encoding/base32"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	cid "github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger"
	crdt "github.com/ipfs/go-ds-crdt"
	"github.com/joho/godotenv"

	"github.com/akakream/sailorsailor/docker"
	"github.com/akakream/sailorsailor/p2p"
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

func (d *Datastore) Shutdown(cancelCtx context.CancelFunc) error {
	cancelCtx()
	log.Println("Shutting down the CRDT")
	err := d.Crdt.Close()
	if err != nil {
		return fmt.Errorf("error while shutting down the CRDT: %w", err)
	}
	log.Println("Shutting down the Store")
	err = d.Store.Close()
	if err != nil {
		return fmt.Errorf("error while shutting down the Datastore: %w", err)
	}
	return nil
}

func (d *Datastore) SetupCRDT(
	ctx context.Context,
	pubsubClient *p2p.LibP2PClient,
	liteipfs *ipfslite.Peer,
) error {
	pubsubBC, err := crdt.NewPubSubBroadcaster(ctx, pubsubClient.Ps, d.topicName)
	if err != nil {
		return err
	}

	opts := crdt.DefaultOptions()
	opts.RebroadcastInterval = 5 * time.Second
	opts.PutHook = putHookLogicCLosure(d, pubsubClient.Host.ID().String())
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

// using docker pull because I need the docker interpretation
func putHookLogicCLosure(d *Datastore, hostID string) func(ds.Key, []byte) {
	return func(k ds.Key, v []byte) {
		fmt.Printf("Added: [%s] -> %s\n", k, string(v))
		if k.Name() == hostID {
			log.Println("HOST IMAGE UPDATE")
			CIDs := strings.Split(string(v), ",")
			cid := CIDs[len(CIDs)-1]
			log.Printf("LAST ADDED CID: %s \n", cid)

			// docker pull tag
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}
			ipdr_server := os.Getenv("IPDR_SERVER")
			dockerizedCid, err := dockerizeCID(cid)
			image_address := fmt.Sprintf("%s/%s", ipdr_server, dockerizedCid)
			fmt.Println(image_address)

			start := time.Now()
			if err := docker.PullCmd(image_address); err != nil {
				fmt.Printf("error while pulling the image: %v\n", err)
			}
			t := time.Now()
			elapsedTime := t.Sub(start)
			fmt.Printf("Docker pull took %s\n", elapsedTime)

			if err := retagImage(d, image_address, cid); err != nil {
				fmt.Printf("error while retagging the image: %v\n", err)
			}
		}
	}
}

func retagImage(d *Datastore, image_address string, cid string) error {
	ctx := context.TODO()
	oldTag := image_address + ":latest"
	k := ds.NewKey(cid)
	v, err := d.Crdt.Get(ctx, k)
	if err != nil {
		return err
	}
	newTag := string(v)
	if err != nil {
		return err
	}
	docker.Tag(oldTag, newTag)
	docker.RemoveTag(oldTag)
	return nil
}

func dockerizeCID(c string) (string, error) {
	decodedCid, err := cid.Decode(c)
	if err != nil {
		return "", err
	}
	decodedHash := []byte(decodedCid.Hash())
	b32str := base32.StdEncoding.EncodeToString(decodedHash)
	end := len(b32str)
	if end > 0 {
		end = end - 1
	}

	// remove padding
	return strings.ToLower(b32str[0:end]), nil
}
