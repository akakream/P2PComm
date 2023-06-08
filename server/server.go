package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/akakream/sailorsailor/identity"
	"github.com/akakream/sailorsailor/p2p"
	store "github.com/akakream/sailorsailor/store"
	ipfslite "github.com/hsanjuan/ipfs-lite"
	"golang.org/x/net/context"
)

func (e apiError) Error() string {
	return e.Err
}

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(apiError); ok {
				writeJSON(w, e.Status, e)
				return
			}
			writeJSON(w, http.StatusInternalServerError, apiError{Err: "internal server error", Status: http.StatusInternalServerError})
		}
	}
}

func NewServer(port string, serverType string, dataPath string, useDatastore bool) *Server {

	ctx, cancel := context.WithCancel(context.Background())

	var servertype ServerType
	var client p2p.P2PClient
	var id *identity.Identity
	var ds *store.Datastore
	var liteIPFS *ipfslite.Peer

	if serverType == "libp2p" {
		servertype = ServerTypeLibp2p
		id, err := identity.NewIdentity(dataPath)
		if err != nil {
			client = p2p.NewLibP2PClient(ctx, false)
		} else {
			client = p2p.NewLibP2PClient(ctx, useDatastore, id)
		}
	} else {
		servertype = ServerTypeIpfs
		client = p2p.NewIpfsP2PClient(&p2p.Config{Port: "5001"})
	}

	if useDatastore {
		if id == nil {
			log.Fatalln("Please provide a data path to use datastore")
		} else {
			d, l, err := setupDataStoreAndIPFSLite(ctx, client, dataPath)
			if err != nil {
				log.Fatal(err)
			}
			ds = d
			liteIPFS = l
		}
	}

	return &Server{
		port:          port,
		Servertype:    servertype,
		DataPath:      dataPath,
		Client:        client,
		Identity:      id,
		Datastore:     ds,
		LiteIPFS:      liteIPFS,
		quitch:        make(chan struct{}),
		cancelContext: cancel,
	}
}

func setupDataStoreAndIPFSLite(ctx context.Context, client p2p.P2PClient, dataPath string) (*store.Datastore, *ipfslite.Peer, error) {
	datastoreTopic := "globaldb-net"
	// Initialize the datastore
	p2pClient, ok := client.(*p2p.LibP2PClient)
	if !ok {
		return nil, nil, errors.New("cannot convert p2p client interface to struct")
	}

	ds, err := store.NewDatastore(ctx, dataPath)
	if err != nil {
		return nil, nil, err
	}
	// Use a special pubsub topic to avoid disconnecting
	// from globaldb peers.
	client.Sub(datastoreTopic)

	liteipfs, err := ipfslite.New(ctx, ds.Store, nil, p2pClient.Host, p2pClient.Dht, nil)
	if err != nil {
		return nil, nil, err
	}

	err = ds.SetupCRDT(ctx, p2pClient, liteipfs)
	if err != nil {
		return nil, nil, err
	}

	return ds, liteipfs, nil
}

func (s *Server) Start() {
	// Publish a message to a topic
	http.HandleFunc("/health", makeHTTPHandler(s.handleHealth))
	// Publish a message to a topic
	http.HandleFunc("/pub", makeHTTPHandler(s.handlePublish))
	// Subscribe to a topic
	http.HandleFunc("/sub", makeHTTPHandler(s.handleSubscribe))
	// Unsubscribe from a topic
	http.HandleFunc("/unsub", makeHTTPHandler(s.handleUnsubscribe))
	// Return all susbcribed topics
	http.HandleFunc("/topics", makeHTTPHandler(s.handleListSubscribedTopics))

	go s.Client.Start()
	go s.listenShutdown()

	go func() {
		if err := http.ListenAndServe(":"+s.port, nil); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	<-s.quitch
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) error {

	var servertype string
	if s.Servertype == ServerTypeLibp2p {
		servertype = "libp2p"
	} else {
		servertype = "ipfs"
	}

	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	return writeJSON(w, http.StatusOK, servertype+" - OK")
}

func (s *Server) gracefullyQuitServer() {
	log.Println("Shutting down the server")

	// Unsub from all topics
	s.Client.Shutdown()
	// Cancel the context
	s.cancelContext()
}

func (s *Server) listenShutdown() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	s.gracefullyQuitServer()
	close(s.quitch)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
