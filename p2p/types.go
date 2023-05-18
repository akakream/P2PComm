package p2p

import (
	shell "github.com/ipfs/go-ipfs-api"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2PClient interface {
	Pub(topicName string, data string)
	Sub(topicName string)
	Unsub(topicName string)
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
}

type Config struct {
	Port string
}
