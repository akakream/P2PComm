package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	ds "github.com/ipfs/go-datastore"
)

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apiError{Err: "invalid body", Status: http.StatusBadRequest}
	}
	defer r.Body.Close()

	var bodyJson Key
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		log.Println(err)
		return apiError{Err: "body must be json", Status: http.StatusBadRequest}
	}

	key := bodyJson.Key
	if key == "" {
		return apiError{Err: "empty key", Status: http.StatusBadRequest}
	}

	ctx := context.TODO()

	// Logic
	value, err := getValue(ctx, s, key)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	resp := KeyValue{
		Key:   key,
		Value: string(value),
	}

	return writeJSON(w, http.StatusOK, resp)
}

func getValue(ctx context.Context, s *Server, key string) ([]byte, error) {
	k := ds.NewKey(key)
	v, err := s.Datastore.Crdt.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[%s] -> %s\n", k, string(v))
	return v, nil
}
