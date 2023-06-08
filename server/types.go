package server

import (
	"context"
	"net/http"

	"github.com/akakream/sailorsailor/identity"
	"github.com/akakream/sailorsailor/p2p"
	store "github.com/akakream/sailorsailor/store"
	ipfslite "github.com/hsanjuan/ipfs-lite"
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
	port          string
	Servertype    ServerType
	DataPath      string
	quitch        chan struct{}
	cancelContext context.CancelFunc
	Client        p2p.P2PClient
	Identity      *identity.Identity
	Datastore     *store.Datastore
	LiteIPFS      *ipfslite.Peer
}

type PubRequestBody struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type PubResponseBody struct {
	Topic    string   `json:"topic"`
	Messages []string `json:"message"`
}

type SubRequestBody struct {
	Topic string `json:"topic"`
}

type UnsubResponseBody struct {
	Topic    string   `json:"topic"`
	Messages []string `json:"messages"`
}

type ListTopicsRequestBody struct {
	Topics []string `json:"topics"`
}

type Key struct {
	Key string `json:"key"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
