package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	ds "github.com/ipfs/go-datastore"
)

func (s *Server) handleCrdtDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := context.TODO()
	key := chi.URLParam(r, "key")

	// Logic
	err := deleteKeyValue(ctx, s, key)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}
	resp := KeyValue{
		Key: key,
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
