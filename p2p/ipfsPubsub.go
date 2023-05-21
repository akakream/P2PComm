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
	topic, topicExists := c.SubscribedTopics[topicName]
	if !topicExists {
		return nil, ErrSubscriptionDoesNotExist
	}

	log.Println("Publishing...")
	if err := c.Shell.PubSubPublish(topicName, data); err != nil {
		log.Printf("Publish failed: %v.\n", err)
		return nil, ErrPublishFailed
	}
	topic.PubHistory = append(topic.PubHistory, data)
	log.Println("Finished publishing.")
	return topic.PubHistory, nil
}

func (c *IpfsP2PClient) Sub(topicName string) error {
	_, topicExists := c.SubscribedTopics[topicName]
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
		c.SubscribedTopics[topicName] = &IpfsP2PTopic{Subscription: subscription, TopicName: topicName, PubHistory: make([]string, 0)}

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

	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		if err := topic.Subscription.Cancel(); err != nil {
			log.Fatal(err)
		}
		delete(c.SubscribedTopics, topicName)
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
		return nil, ErrSubscriptionDoesNotExist
	}
	return topic.PubHistory, nil
}

func (c *IpfsP2PClient) Shutdown() {
	for _, topic := range c.SubscribedTopics {
		c.Unsub(topic.TopicName)
	}
	// close(c.Channel)
	// log.Println("Closing channel.")
	time.Sleep(2 * time.Second)
}

func (c *IpfsP2PClient) ListSubscribedTopics() []string {
	subscribedTopics := make([]string, 0)
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.TopicName)
	}
	return subscribedTopics
}
