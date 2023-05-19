package p2p

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (c *LibP2PClient) listen(channel chan *pubsub.Message) {
	for {
		select {
		case msg := <-channel:
			fmt.Println(string(msg.Data))
			fmt.Println(msg.Topic)
		default:
		}
	}
}

func NewLibP2PClient() *LibP2PClient {

	ctx := context.Background()
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

	return &LibP2PClient{
		Host:             host,
		Ps:               ps,
		SubscribedTopics: make(map[string]LibP2PTopic),
		Channel:          make(chan *pubsub.Message),
	}
}

func (c *LibP2PClient) Start() {
	c.listen(c.Channel)
}

func (c *LibP2PClient) Pub(topicName string, data string) {
	ctx := context.Background()
	log.Println("Publishing...")
	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		if err := topic.Topic.Publish(ctx, []byte(data)); err != nil {
			log.Fatal(err)
		}
		log.Println("Finished publishing.")
	}
}

func (c *LibP2PClient) Sub(topicName string, wg *sync.WaitGroup) {
	ctx := context.Background()
	topic, err := c.Ps.Join(topicName)
	if err != nil {
		log.Fatal(err)
	}

	subscription, err := topic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}

	c.SubscribedTopics[topicName] = LibP2PTopic{Subscription: subscription, Topic: topic}

	wg.Done()

	for {
		msg, err := subscription.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		c.Channel <- msg
	}
}

func (c *LibP2PClient) Unsub(topicName string) {

	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		topic.Subscription.Cancel()
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
	}
}

func (c *LibP2PClient) Shutdown() {
	for _, topic := range c.SubscribedTopics {
		c.Unsub(topic.Topic.String())
	}
}

func (c *LibP2PClient) ListSubscribedTopics() []string {
	var subscribedTopics []string
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.Topic.String())
	}
	return subscribedTopics
}
