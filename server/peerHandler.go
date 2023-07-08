package server

import (
	"net/http"
)

// Get the key stored in the datastore
func (s *Server) handlePeersGet(w http.ResponseWriter, r *http.Request) error {
	// ctx := context.TODO()

	// Logic
	ps := s.Client.Host.Peerstore()
	peers := ps.Peers()

	/*
		if err != nil {
			return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
		}
	*/

	resp := struct {
		Peers string `json:"peers"`
	}{
		Peers: peers.String(),
	}

	return writeJSON(w, http.StatusOK, resp)
}
