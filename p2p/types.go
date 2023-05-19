package p2p

import (
	"sync"

	shell "github.com/ipfs/go-ipfs-api"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2PClient interface {
	Start()
	Pub(topicName string, data string)
	Sub(topicName string, wg *sync.WaitGroup)
	Unsub(topicName string)
	Shutdown()
	ListSubscribedTopics() []string
}

type LibP2PClient struct {
	Host             host.Host
	Ps               *pubsub.PubSub
	SubscribedTopics map[string]LibP2PTopic
	Channel          chan *pubsub.Message
}

type LibP2PTopic struct {
	Subscription *pubsub.Subscription
	Topic        *pubsub.Topic
}

type IpfsP2PClient struct {
	Shell            *shell.Shell
	SubscribedTopics map[string]IpfsP2PTopic
	Channel          chan *shell.Message
}

type IpfsP2PTopic struct {
	Subscription *shell.PubSubSubscription
	TopicName    string
}

type Config struct {
	Port string
}
