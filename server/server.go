package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/akakream/sailorsailor/identity"
	"github.com/akakream/sailorsailor/p2p"
	store "github.com/akakream/sailorsailor/store"
	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func (e apiError) Error() string {
	return e.Err
}

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(apiError); ok {
				err := writeJSON(w, e.Status, e)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			err := writeJSON(w, http.StatusInternalServerError, apiError{Err: "internal server error", Status: http.StatusInternalServerError})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func NewServer(port string, dataPath string, useDatastore bool) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	var client *p2p.LibP2PClient
	var id *identity.Identity
	var ds *store.Datastore
	var liteIPFS *ipfslite.Peer

	identity, err := identity.NewIdentity(dataPath)
	if err != nil {
		client = p2p.NewLibP2PClient(ctx, false)
	} else {
		client = p2p.NewLibP2PClient(ctx, useDatastore, identity)
		id = identity
	}

	if useDatastore {
		if id == nil {
			log.Fatalln("Please provide a data path to use datastore")
		} else {
			d, l, err := setupDataStoreAndIPFSLite(ctx, client, filepath.Join(dataPath, "datastore"))
			if err != nil {
				log.Fatal(err)
			}
			ds = d
			liteIPFS = l
		}
	}

	return &Server{
		port:          port,
		DataPath:      dataPath,
		Client:        client,
		Identity:      id,
		Datastore:     ds,
		LiteIPFS:      liteIPFS,
		quitch:        make(chan struct{}),
		cancelContext: cancel,
	}
}

func setupDataStoreAndIPFSLite(ctx context.Context, client *p2p.LibP2PClient, dataPath string) (*store.Datastore, *ipfslite.Peer, error) {
	datastoreTopic := "globaldb-net"
	// Initialize the datastore

	ds, err := store.NewDatastore(ctx, dataPath)
	if err != nil {
		return nil, nil, err
	}
	// Use a special pubsub topic to avoid disconnecting
	// from globaldb peers.
	err = client.Sub(datastoreTopic)
	if err != nil {
		return nil, nil, err
	}

	liteipfs, err := ipfslite.New(ctx, ds.Store, nil, client.Host, client.Dht, nil)
	if err != nil {
		return nil, nil, err
	}

	err = ds.SetupCRDT(ctx, client, liteipfs)
	if err != nil {
		return nil, nil, err
	}

    // Ping every 20 seconds to keep the topic alive
    go ping(ctx, client, datastoreTopic)
    // client.Host.ConnManager().Protect("peerid", "keep")

    bootstrap(client, liteipfs)

	return ds, liteipfs, nil
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// Publish a message to a topic
	r.Get("/health", makeHTTPHandler(s.handleHealth))
	// Publish a message to a topic
	r.Post("/pub", makeHTTPHandler(s.handlePublish))
	// Subscribe to a topic
	r.Post("/sub", makeHTTPHandler(s.handleSubscribe))
	// Unsubscribe from a topic
	r.Post("/unsub", makeHTTPHandler(s.handleUnsubscribe))
	// Return all susbcribed topics
	r.Get("/topics", makeHTTPHandler(s.handleListSubscribedTopics))

	r.Get("/crdt", makeHTTPHandler(s.handleCrdtGet))
	r.Get("/crdt/{key}", makeHTTPHandler(s.handleGetCrdtByID))
	r.Post("/crdt", makeHTTPHandler(s.handleCrdtPost))
	r.Delete("/crdt/{key}", makeHTTPHandler(s.handleCrdtDelete))

	r.Get("/peers", makeHTTPHandler(s.handlePeersGet))
	r.Get("/identity", makeHTTPHandler(s.handleIdentityGet))

	go s.Client.Start()
	go s.listenShutdown()

	go func() {
		if err := http.ListenAndServe(":"+s.port, r); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	<-s.quitch
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	return writeJSON(w, http.StatusOK, "OK")
}

func (s *Server) gracefullyQuitServer() {
	log.Println("Shutting down the server")

	// Shutdown the datastore
	if s.Datastore != nil {
		err := s.Datastore.Shutdown()
		if err != nil {
			err = fmt.Errorf("error while shutting down the server gracefully: %w", err)
			log.Fatal(err)
		}
	}
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

func bootstrap(client *p2p.LibP2PClient, liteipfs *ipfslite.Peer) {
	log.Println("Bootstrapping...")
	// peersList := ipfslite.DefaultBootstrapPeers()
	peersList := []peer.AddrInfo{}

	// Read peers from the peerstore
	peersFile, err := os.OpenFile("./data/peerstore", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("reading peerstore file is unsuccessfull. Add at least one peer to connect with others. Error: %v", err)
	}
	defer peersFile.Close()

	scanner := bufio.NewScanner(peersFile)
	for scanner.Scan() {
		peerAddress := scanner.Text()
		log.Printf("Adding peer %s to peer list...", peerAddress)
		bstr, _ := multiaddr.NewMultiaddr(peerAddress)
		inf, _ := peer.AddrInfoFromP2pAddr(bstr)
		peersList = append(peersList, *inf)
		client.Host.ConnManager().TagPeer(inf.ID, "keep", 100)
	}
	log.Println("Bootstrapping following peers: ", peersList)

	liteipfs.Bootstrap(peersList)
	log.Println("Bootstrapping done.")
}

func ping(ctx context.Context, client *p2p.LibP2PClient, topic string) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            client.Pub(topic, "ping")
            time.Sleep(20 * time.Second)
        }
    }
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
