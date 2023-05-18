package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akakream/sailorsailor/p2p"
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

func NewServer(port string, serverType string) *Server {

	// ctx := context.Background()

	var servertype ServerType
	var client p2p.P2PClient
	if serverType == "libp2p" {
		servertype = ServerTypeLibp2p
		client = p2p.NewLibP2PClient()
	} else {
		servertype = ServerTypeIpfs
		client = p2p.NewIpfsP2PClient(&p2p.Config{Port: "5001"})
	}

	return &Server{
		port:       port,
		Servertype: servertype,
		Client:     client,
		quitch:     make(chan struct{}),
	}
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

	http.ListenAndServe(":"+s.port, nil)

	<-s.quitch
	fmt.Println("Shutting down the server")
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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
