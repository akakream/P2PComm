package p2p

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/akakream/sailorsailor/identity"
	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
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

func NewLibP2PClient(ctx context.Context, liteipfs bool, id ...*identity.Identity) *LibP2PClient {
	var host host.Host
	var dht *dual.DHT
	listen, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/33123")
	if err != nil {
		log.Fatal(err)
	}

	if liteipfs {
		h, d, err := ipfslite.SetupLibp2p(ctx, id[0].PrivKey, nil, []multiaddr.Multiaddr{listen}, nil, ipfslite.Libp2pOptionsExtra...)
		if err != nil {
			log.Fatal(err)
		}
		host = h
		dht = d
	} else {
		if len(id) == 0 {
			h, err := libp2p.New(libp2p.ListenAddrStrings(listen.String()))
			if err != nil {
				log.Fatal(err)
			}
			host = h
		} else {
			h, err := libp2p.New(libp2p.Identity(id[0].PrivKey), libp2p.ListenAddrStrings(listen.String()))
			if err != nil {
				log.Fatal(err)
			}
			host = h
		}
	}
	log.Println("HOST ID: ", host.ID())
	log.Println("HOST ADDR: ", host.Addrs())

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	return &LibP2PClient{
		Host:             host,
		Dht:              dht,
		Ps:               ps,
		SubscribedTopics: make(map[string]*LibP2PTopic),
		Channel:          make(chan *pubsub.Message, 20),
	}
}

func (c *LibP2PClient) Start() {
	c.listen(c.Channel)
}

func (c *LibP2PClient) Pub(topicName string, data string) ([]string, error) {
	c.mu.RLock()
	topic, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()
	if !topicExists {
		return nil, ErrSubscriptionDoesNotExist
	}

	log.Println("Publishing...")
	ctx := context.Background()
	if err := topic.Topic.Publish(ctx, []byte(data)); err != nil {
		log.Printf("Publish failed: %v.\n", err)
		return nil, ErrPublishFailed
	}
	topic.mu.Lock()
	topic.PubHistory = append(topic.PubHistory, data)
	topic.mu.Unlock()
	log.Println("Finished publishing.")
	return topic.PubHistory, nil
}

func (c *LibP2PClient) Sub(topicName string) error {
	c.mu.RLock()
	_, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()
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

		c.mu.Lock()
		c.SubscribedTopics[topicName] = &LibP2PTopic{Subscription: subscription, Topic: topic, PubHistory: make([]string, 0)}
		c.mu.Unlock()

		wg.Done()

		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				fmt.Printf("fetching subscription message error: %v.\n", err)
				break
			}
			c.Host.ConnManager().TagPeer(msg.ReceivedFrom, "keep", 100)
			c.Channel <- msg
		}
	}()
	wg.Wait()
	return nil
}

func (c *LibP2PClient) Unsub(topicName string) ([]string, error) {
	c.mu.RLock()
	topic, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()

	if topicExists {
		topic.Subscription.Cancel()
		topic.Topic.Close()
		// TODO: THIS BLOCKS, I DONT KNOW WHY???
		//c.mu.Lock()
		delete(c.SubscribedTopics, topicName)
		//c.mu.Unlock()
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
		return nil, ErrSubscriptionDoesNotExist
	}

	return topic.PubHistory, nil
}

func (c *LibP2PClient) Shutdown() {
	if c.Dht != nil {
		log.Println("Closing DHT...")
		c.Dht.Close()
	}

	c.mu.RLock()
	log.Println("Unsubscribing from the topics...")
	for _, topic := range c.SubscribedTopics {
		_, err := c.Unsub(topic.Topic.String())
		if err != nil {
			log.Printf("Could not unsubscribe from %s", topic.Topic.String())
		}
	}
	c.mu.RUnlock()

	// close(c.Channel)
	// log.Println("Closing channel.")
	log.Println("Closing host...")
	c.Host.Close()
	time.Sleep(2 * time.Second)
}

func (c *LibP2PClient) ListSubscribedTopics() []string {
	subscribedTopics := make([]string, 0)
	c.mu.RLock()
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.Topic.String())
	}
	c.mu.RUnlock()
	return subscribedTopics
}
