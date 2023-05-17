package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiError struct {
	Err    string `json:"err"`
	Status int    `json:"status"`
}

func (e apiError) Error() string {
	return e.Err
}

type apiFunc func(http.ResponseWriter, *http.Request) error

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

type Server struct {
	port   string
	quitch chan struct{}
}

func NewServer(port string) *Server {
	return &Server{
		port:   port,
		quitch: make(chan struct{}),
	}
}

func (s *Server) Start() {
	// Publish a message to a topic
	http.HandleFunc("/health", makeHTTPHandler(handleHealth))
	// Publish a message to a topic
	http.HandleFunc("/pub", makeHTTPHandler(handlePublish))
	// Subscribe to a topic
	http.HandleFunc("/sub", makeHTTPHandler(handleSubscribe))
	// Unsubscribe from a topic
	http.HandleFunc("/unsub", makeHTTPHandler(handleUnsubscribe))

	http.ListenAndServe(":"+s.port, nil)

	<-s.quitch
	fmt.Println("Shutting down the server")
}

func handlePublish(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func handleUnsubscribe(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	return writeJSON(w, http.StatusOK, "OK")
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
