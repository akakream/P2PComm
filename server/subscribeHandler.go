package server

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) error {
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
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go s.Client.Sub(topicName, wg)
	wg.Wait()

	return writeJSON(w, http.StatusOK, "Subscribed to the topic "+bodyJson.Topic)
}
