package server

import (
	"net/http"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func (s *Server) handlePeersGet(w http.ResponseWriter, r *http.Request) error {
	// ctx := context.TODO()

	// Logic
	// ps := s.Client.Host.Peerstore()
    // peers := ps.Peers()
    //
    /*
    peers := s.Client.Ps.ListPeers("globaldb-net")
    peers = append(peers, s.Client.Host.ID())
    var peerIDs []string
    for _, peer := range peers {
        peerIDs = append(peerIDs, peer.String())
    }
    */
    peers := connectedPeers(s.Client.Host)
    peers = append(peers, &peer.AddrInfo{
        ID:    s.Client.Host.ID(),
        Addrs: s.Client.Host.Addrs(),
    })
	/*
		if err != nil {
			return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
		}
	*/

	resp := struct {
		Peers []*peer.AddrInfo `json:"peers"`
	}{
		Peers: peers,
	}

	return writeJSON(w, http.StatusOK, resp)
}


func connectedPeers(h host.Host) []*peer.AddrInfo {
	var pinfos []*peer.AddrInfo
	for _, c := range h.Network().Conns() {
		pinfos = append(pinfos, &peer.AddrInfo{
			ID:    c.RemotePeer(),
			Addrs: []multiaddr.Multiaddr{c.RemoteMultiaddr()},
		})
	}
	return pinfos
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
