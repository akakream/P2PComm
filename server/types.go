package server

import (
	"net/http"

	"github.com/akakream/sailorsailor/p2p"
)

type apiError struct {
	Err    string `json:"err"`
	Status int    `json:"status"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ServerType byte

const (
	ServerTypeIpfs ServerType = iota
	ServerTypeLibp2p
)

type Server struct {
	port       string
	Servertype ServerType
	quitch     chan struct{}
	client     *p2p.P2PClient
}

type PubRequestBody struct {
	Topic   string `json:"topic"`
	Message []byte `json:"message"`
}

type SubRequestBody struct {
	Topic string `json:"topic"`
}
