package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	ds "github.com/ipfs/go-datastore"

	"github.com/akakream/sailorsailor/utils"
)

// Get the key stored in the datastore
func (s *Server) handleGetCrdtByID(w http.ResponseWriter, r *http.Request) error {
	ctx := context.TODO()
	key := chi.URLParam(r, "key")
	decodedKey, err := utils.DecodeParam(key)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	// Logic
	value, err := GetValue(ctx, s, decodedKey)
	if err != nil {
		return apiError{Err: "Key not found", Status: http.StatusInternalServerError}
	}

	resp := KeyValue{
		Key:   decodedKey,
		Value: string(value),
	}

	return writeJSON(w, http.StatusOK, resp)
}

func GetValue(ctx context.Context, s *Server, key string) ([]byte, error) {
	k := ds.NewKey(key)
	v, err := s.Datastore.Crdt.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[%s] -> %s\n", k, string(v))
	return v, nil
}
