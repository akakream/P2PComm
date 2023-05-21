package p2p

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (c *LibP2PClient) listen(channel chan *pubsub.Message) {
	for {
		select {
		case msg := <-channel:
			log.Println(string(msg.Data))
			log.Println(msg.Topic)
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
	log.Println(host.ID())
	log.Println(host.Addrs())

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	return &LibP2PClient{
		Host:             host,
		Ps:               ps,
		SubscribedTopics: make(map[string]*LibP2PTopic),
		Channel:          make(chan *pubsub.Message, 20),
	}
}

func (c *LibP2PClient) Start() {
	c.listen(c.Channel)
}

func (c *LibP2PClient) Pub(topicName string, data string) ([]string, error) {
	topic, topicExists := c.SubscribedTopics[topicName]
	if !topicExists {
		return nil, ErrSubscriptionDoesNotExist
	}

	log.Println("Publishing...")
	ctx := context.Background()
	if err := topic.Topic.Publish(ctx, []byte(data)); err != nil {
		log.Printf("Publish failed: %v.\n", err)
		return nil, ErrPublishFailed
	}
	topic.PubHistory = append(topic.PubHistory, data)
	log.Println("Finished publishing.")
	return topic.PubHistory, nil
}

func (c *LibP2PClient) Sub(topicName string) error {
	_, topicExists := c.SubscribedTopics[topicName]
	if topicExists {
		return ErrAlreadySubscribed
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		ctx := context.Background()
		topic, err := c.Ps.Join(topicName)
		if err != nil {
			log.Fatalf("join error: %v", err)
		}

		subscription, err := topic.Subscribe()
		if err != nil {
			log.Fatalf("subscription error: %v", err)
		}

		c.SubscribedTopics[topicName] = &LibP2PTopic{Subscription: subscription, Topic: topic, PubHistory: make([]string, 0)}

		wg.Done()

		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				fmt.Printf("fetching subscription message error: %v.\n", err)
				break
			}
			c.Channel <- msg
		}
	}()
	wg.Wait()
	return nil
}

func (c *LibP2PClient) Unsub(topicName string) ([]string, error) {
	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		topic.Subscription.Cancel()
		topic.Topic.Close()
		delete(c.SubscribedTopics, topicName)
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
		return nil, ErrSubscriptionDoesNotExist
	}

	return topic.PubHistory, nil
}

func (c *LibP2PClient) Shutdown() {
	for _, topic := range c.SubscribedTopics {
		c.Unsub(topic.Topic.String())
	}
	// close(c.Channel)
	// log.Println("Closing channel.")
	c.Host.Close()
	log.Println("Closing host.")
	time.Sleep(2 * time.Second)
}

func (c *LibP2PClient) ListSubscribedTopics() []string {
	subscribedTopics := make([]string, 0)
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.Topic.String())
	}
	return subscribedTopics
}
