package libp2pPubsub

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type PubSubClient struct {
	Host             host.Host
	Ps               *pubsub.PubSub
	SubscribedTopics map[string]Topic
	Channel          chan *pubsub.Message
}

type Topic struct {
	Subscription *pubsub.Subscription
	Topic        *pubsub.Topic
}

func Listen(channel chan *pubsub.Message) {
	for {
		select {
		case msg := <-channel:
			fmt.Println(string(msg.Data))
			fmt.Println(msg.Topic)
		default:
		}
	}
}

func Run() {
	ctx := context.Background()
	topicName := "babuska1"
	pubsubClient := NewPubSubClient(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go pubsubClient.Sub(topicName, ctx, wg)
	time.Sleep(5 * time.Second)
	go Listen(pubsubClient.Channel)

	go pubsubClient.Pub(ctx, topicName, "lolwut1")
	go pubsubClient.Pub(ctx, topicName, "lolwut2")
	go pubsubClient.Pub(ctx, topicName, "lolwut3")
	go pubsubClient.Pub(ctx, topicName, "lolwut4")
	go pubsubClient.Pub(ctx, topicName, "lolwut5")
	go pubsubClient.Pub(ctx, topicName, "lolwut6")

	wg.Wait()
	pubsubClient.Unsub(topicName)

}

func NewPubSubClient(ctx context.Context) *PubSubClient {

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	fmt.Println(host.ID())
	fmt.Println(host.Addrs())

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	return &PubSubClient{
		Host:             host,
		Ps:               ps,
		SubscribedTopics: make(map[string]Topic),
		Channel:          make(chan *pubsub.Message),
	}
}

func (c *PubSubClient) Pub(ctx context.Context, topicName string, data string) {
	log.Println("Publishing...")
	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		if err := topic.Topic.Publish(ctx, []byte(data)); err != nil {
			log.Fatal(err)
		}
		log.Println("Finished publishing.")
	}
}

func (c *PubSubClient) Sub(topicName string, ctx context.Context, wg *sync.WaitGroup) {
	topic, err := c.Ps.Join(topicName)
	if err != nil {
		log.Fatal(err)
	}

	subscription, err := topic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}

	c.SubscribedTopics[topicName] = Topic{Subscription: subscription, Topic: topic}

	for {
		msg, err := subscription.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		c.Channel <- msg
	}
	wg.Done()
}

func (c *PubSubClient) Unsub(topicName string) {

	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		topic.Subscription.Cancel()
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
	}
}
