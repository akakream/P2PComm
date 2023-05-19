package p2p

import (
	"fmt"
	"log"
	"sync"

	shell "github.com/ipfs/go-ipfs-api"
)

func (c *IpfsP2PClient) listen(channel chan *shell.Message) {
	for {
		select {
		case msg := <-channel:
			fmt.Println(string(msg.Data))
			fmt.Println(msg.TopicIDs)
		default:
		}
	}
}

func NewIpfsP2PClient(config *Config) *IpfsP2PClient {
	// ipfs daemon --enable-pubsub-experiment
	sh := shell.NewShell("localhost:" + config.Port)
	return &IpfsP2PClient{
		Shell:            sh,
		SubscribedTopics: make(map[string]IpfsP2PTopic),
		Channel:          make(chan *shell.Message),
	}
}

func (c *IpfsP2PClient) Start() {
	c.listen(c.Channel)
}

func (c *IpfsP2PClient) Pub(topicName string, data string) {
	log.Println("Publishing...")
	if err := c.Shell.PubSubPublish(topicName, data); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished publishing.")
}

func (c *IpfsP2PClient) Sub(topicName string, wg *sync.WaitGroup) {
	subscription, err := c.Shell.PubSubSubscribe(topicName)
	if err != nil {
		log.Fatal(err)
	}
	c.SubscribedTopics[topicName] = IpfsP2PTopic{Subscription: subscription, TopicName: topicName}

	wg.Done()

	for {
		msg, err := subscription.Next()
		if err != nil {
			log.Fatal(err)
		}
		c.Channel <- msg
	}
}

func (c *IpfsP2PClient) Unsub(topicName string) {

	topic, topicExists := c.SubscribedTopics[topicName]

	if topicExists {
		if err := topic.Subscription.Cancel(); err != nil {
			log.Fatal(err)
		}
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
	}
}

func (c *IpfsP2PClient) Shutdown() {
	for _, topic := range c.SubscribedTopics {
		c.Unsub(topic.TopicName)
	}
}

func (c *IpfsP2PClient) ListSubscribedTopics() []string {
	var subscribedTopics []string
	for _, topic := range c.SubscribedTopics {
		subscribedTopics = append(subscribedTopics, topic.TopicName)
	}
	return subscribedTopics
}
