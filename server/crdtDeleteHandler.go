package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	ds "github.com/ipfs/go-datastore"

	"github.com/akakream/sailorsailor/utils"
)

func (s *Server) handleCrdtDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := context.TODO()
	key := chi.URLParam(r, "key")
	decodedKey, err := utils.DecodeParam(key)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	// Logic
	err = deleteKeyValue(ctx, s, decodedKey)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}
	resp := KeyValue{
		Key: decodedKey,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func deleteKeyValue(ctx context.Context, s *Server, key string) error {
	k := ds.NewKey(key)
	err := s.Datastore.Crdt.Delete(ctx, k)
	if err != nil {
		return err
	}
	return nil
}
