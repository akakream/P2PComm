package p2p

import (
	shell "github.com/ipfs/go-ipfs-api"
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
	Ps               *pubsub.PubSub
	SubscribedTopics map[string]*LibP2PTopic
	Channel          chan *pubsub.Message
}

type LibP2PTopic struct {
	Subscription *pubsub.Subscription
	Topic        *pubsub.Topic
	PubHistory   []string
}

type IpfsP2PClient struct {
	Shell            *shell.Shell
	SubscribedTopics map[string]*IpfsP2PTopic
	Channel          chan *shell.Message
}

type IpfsP2PTopic struct {
	Subscription *shell.PubSubSubscription
	TopicName    string
	PubHistory   []string
}

type Config struct {
	Port string
}
