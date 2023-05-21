package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) handleUnsubscribe(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apiError{Err: "invalid body", Status: http.StatusBadRequest}
	}
	defer r.Body.Close()

	var bodyJson SubRequestBody
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return apiError{Err: "body must be json", Status: http.StatusBadRequest}
	}

	topicName := bodyJson.Topic
	if topicName == "" {
		return apiError{Err: "empty topic", Status: http.StatusBadRequest}
	}

	// Logic
	messageHistory, err := s.Client.Unsub(topicName)
	if err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}
	resp := UnsubResponseBody{
		Topic:    topicName,
		Messages: messageHistory,
	}

	return writeJSON(w, http.StatusOK, resp)
}
