package p2p

import (
	"sync"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2PClient interface {
	Start()
	Pub(topicName string, data string) ([]string, error)
	Sub(topicName string) error
	Unsub(topicName string) ([]string, error)
	Shutdown()
	ListSubscribedTopics() []string
}

type LibP2PClient struct {
	Host             host.Host
	Dht              *dual.DHT
	Ps               *pubsub.PubSub
	mu               sync.RWMutex
	SubscribedTopics map[string]*LibP2PTopic
	Channel          chan *pubsub.Message
}

type LibP2PTopic struct {
	Subscription *pubsub.Subscription
	Topic        *pubsub.Topic
	mu           sync.RWMutex
	PubHistory   []string
}

type IpfsP2PClient struct {
	Shell            *shell.Shell
	mu               sync.RWMutex
	SubscribedTopics map[string]*IpfsP2PTopic
	Channel          chan *shell.Message
}

type IpfsP2PTopic struct {
	Subscription *shell.PubSubSubscription
	TopicName    string
	mu           sync.RWMutex
	PubHistory   []string
}

type Config struct {
	Port string
}
