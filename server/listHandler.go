package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ipfs/go-datastore/query"
)

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	ctx := context.TODO()

	// Logic
	resp, err := getKeyValues(ctx, s)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return writeJSON(w, http.StatusOK, *resp)
}

func getKeyValues(ctx context.Context, s *Server) (*[]KeyValue, error) {
	q := query.Query{}
	results, err := s.Datastore.Crdt.Query(ctx, q)
	if err != nil {
		return nil, apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}
	result := []KeyValue{}
	for r := range results.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		pair := KeyValue{
			Key:   r.Key,
			Value: string(r.Value),
		}
		result = append(result, pair)

		fmt.Printf("[%s] -> %s\n", r.Key, string(r.Value))
	}

	return &result, nil
}
