package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	ds "github.com/ipfs/go-datastore"
)

func (s *Server) handleCrdtPost(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apiError{Err: "invalid body", Status: http.StatusBadRequest}
	}
	defer r.Body.Close()

	var bodyJson KeyValue
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		log.Println(err)
		return apiError{Err: "body must be json", Status: http.StatusBadRequest}
	}

	key := bodyJson.Key
	if key == "" {
		return apiError{Err: "empty key", Status: http.StatusBadRequest}
	}

	value := bodyJson.Value
	if value == "" {
		return apiError{Err: "empty value", Status: http.StatusBadRequest}
	}

	ctx := context.TODO()

	// Logic
	err = putKeyValue(ctx, s, key, string(value))
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}
	resp := KeyValue{
		Key:   key,
		Value: value,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func putKeyValue(ctx context.Context, s *Server, key string, value string) error {
	k := ds.NewKey(key)
	err := s.Datastore.Crdt.Put(ctx, k, []byte(value))
	if err != nil {
		return err
	}
	return nil
}
