package server

import (
	"net/http"
)

func (s *Server) handleListSubscribedTopics(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	// Logic
	topics := s.Client.ListSubscribedTopics()
	resp := ListTopicsRequestBody{
		Topics: topics,
	}

	return writeJSON(w, http.StatusOK, resp)
}
