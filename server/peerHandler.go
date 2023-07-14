package server

import (
	"net/http"
)

func (s *Server) handlePeersGet(w http.ResponseWriter, r *http.Request) error {
	// ctx := context.TODO()

	// Logic
	// ps := s.Client.Host.Peerstore()
    // peers := ps.Peers()
    peers := s.Client.Ps.ListPeers("globaldb-net")
    var peerIDs []string
    for _, peer := range peers {
        peerIDs = append(peerIDs, peer.String())
    }

	/*
		if err != nil {
			return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
		}
	*/

	resp := struct {
		Peers []string `json:"peers"`
	}{
		Peers: peerIDs,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleIdentityGet(w http.ResponseWriter, r *http.Request) error {
	// ctx := context.TODO()

	// Logic
    id := s.Client.Host.ID()
    addrs := s.Client.Host.Addrs()
    var hostAddresses []string
    for _, addr := range addrs {
        hostAddresses = append(hostAddresses, addr.String())
    }

	resp := struct {
		ID string `json:"id"`
		Addrs []string `json:"addrs"`
	}{
        ID: id.String(),
        Addrs: hostAddresses,
	}

	return writeJSON(w, http.StatusOK, resp)
}
