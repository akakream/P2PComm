package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	if bodyJson.Topic == "" {
		return apiError{Err: "empty topic", Status: http.StatusBadRequest}
	}

	// Logic
	if s.Servertype == ServerTypeLibp2p {
		fmt.Println("libp2p")
	} else {
		fmt.Println("ipfs")
	}

	return writeJSON(w, http.StatusOK, "Subscribed to the topic "+bodyJson.Topic)
}
