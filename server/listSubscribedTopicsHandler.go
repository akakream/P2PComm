package server

import (
	"net/http"
)

type response struct {
	Topics []string `json:"topics"`
}

func (s *Server) handleListSubscribedTopics(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	// Logic
	topics := s.Client.ListSubscribedTopics()
	resp := response{
		Topics: topics,
	}

	return writeJSON(w, http.StatusOK, resp)
}
