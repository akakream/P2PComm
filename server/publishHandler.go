package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apiError{Err: "invalid body", Status: http.StatusBadRequest}
	}
	defer r.Body.Close()

	var bodyJson PubRequestBody
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		log.Println(err)
		return apiError{Err: "body must be json", Status: http.StatusBadRequest}
	}

	topicName := bodyJson.Topic
	if topicName == "" {
		return apiError{Err: "empty topic", Status: http.StatusBadRequest}
	}

	message := bodyJson.Message
	if message == "" {
		return apiError{Err: "empty message", Status: http.StatusBadRequest}
	}

	// Logic
	if err := s.Client.Pub(topicName, message); err != nil {
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return writeJSON(w, http.StatusOK, "Published message "+string(message)+" to the topic "+bodyJson.Topic)
}
