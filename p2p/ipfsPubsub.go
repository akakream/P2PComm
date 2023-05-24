package p2p

import (
	"log"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func (c *IpfsP2PClient) listen(channel chan *shell.Message) {
	for {
		select {
		case msg := <-channel:
			log.Println(string(msg.Data))
			log.Println(msg.TopicIDs)
		default:
		}
	}
}

func NewIpfsP2PClient(config *Config) *IpfsP2PClient {
	// ipfs daemon --enable-pubsub-experiment
	sh := shell.NewShell("localhost:" + config.Port)
	return &IpfsP2PClient{
		Shell:            sh,
		SubscribedTopics: make(map[string]*IpfsP2PTopic),
		Channel:          make(chan *shell.Message, 20),
	}
}

func (c *IpfsP2PClient) Start() {
	c.listen(c.Channel)
}

func (c *IpfsP2PClient) Pub(topicName string, data string) ([]string, error) {
	c.mu.RLock()
	topic, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()
	if !topicExists {
		return nil, ErrSubscriptionDoesNotExist
	}

	log.Println("Publishing...")
	if err := c.Shell.PubSubPublish(topicName, data); err != nil {
		log.Printf("Publish failed: %v.\n", err)
		return nil, ErrPublishFailed
	}
	topic.mu.Lock()
	topic.PubHistory = append(topic.PubHistory, data)
	topic.mu.Unlock()
	log.Println("Finished publishing.")
	return topic.PubHistory, nil
}

func (c *IpfsP2PClient) Sub(topicName string) error {
	c.mu.RLock()
	_, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()
	if topicExists {
		return ErrAlreadySubscribed
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		subscription, err := c.Shell.PubSubSubscribe(topicName)
		if err != nil {
			log.Fatal(err)
		}
		c.mu.Lock()
		c.SubscribedTopics[topicName] = &IpfsP2PTopic{Subscription: subscription, TopicName: topicName, PubHistory: make([]string, 0)}
		c.mu.Unlock()

		wg.Done()

		for {
			msg, err := subscription.Next()
			if err != nil {
				log.Printf("fetching subscription message error: %v.\n", err)
				break
			}
			c.Channel <- msg
		}
	}()
	wg.Wait()
	return nil
}

func (c *IpfsP2PClient) Unsub(topicName string) ([]string, error) {
	c.mu.RLock()
	topic, topicExists := c.SubscribedTopics[topicName]
	c.mu.RUnlock()

	if topicExists {
		if err := topic.Subscription.Cancel(); err != nil {
			log.Fatal(err)
		}
		c.mu.Lock()
		delete(c.SubscribedTopics, topicName)
		c.mu.Unlock()
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
		return nil, ErrSubscriptionDoesNotExist
	}
	return topic.PubHistory, nil
}

func (c *IpfsP2PClient) Shutdown() {
	c.mu.RLock()
	for _, topic := range c.SubscribedTopics {
		c.Unsub(topic.TopicName)
	}
	c.mu.RUnlock()
	// close(c.Channel)
	// log.Println("Closing channel.")
	time.Sleep(2 * time.Second)
}

func (c *IpfsP2PClient) ListSubscribedTopics() []string {
	subscribedTopics := make([]string, 0)
	c.mu.RLock()
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.TopicName)
	}
	c.mu.RUnlock()
	return subscribedTopics
}
