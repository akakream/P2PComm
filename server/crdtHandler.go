package server

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/ipfs/go-datastore/query"
)

func (s *Server) handleCrdtGet(w http.ResponseWriter, r *http.Request) error {
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
        key := trimTheSlashInTheBeginning(r.Key)
		pair := KeyValue{
			Key:   key,
			Value: string(r.Value),
		}
		result = append(result, pair)

		fmt.Printf("[%s] -> %s\n", key, string(r.Value))
	}

	return &result, nil
}

func trimTheSlashInTheBeginning(key string) string {
        c, i := utf8.DecodeRuneInString(key)
        if c == []rune("/")[0] {
            return key[i:]
        }
        return key
}
