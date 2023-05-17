package ipfsPubsub

import (
	"fmt"
	"log"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type PubSubClient struct {
	Shell            *shell.Shell
	SubscribedTopics map[string]Topic
	Channel          chan *shell.Message
}

type Config struct {
	Port string
}

type Topic struct {
	Subscription *shell.PubSubSubscription
}

func Listen(channel chan *shell.Message) {
	for {
		select {
		case msg := <-channel:
			fmt.Println(string(msg.Data))
			fmt.Println(msg.TopicIDs)
		default:
		}
	}
}

func Run() {
	sh := shell.NewShell("localhost:" + "5001")
	_ = sh

	topicName := "babuska1"
	pubsubClient := NewPubSubClient(&Config{Port: "5001"})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go pubsubClient.Sub(topicName, wg)
	time.Sleep(5 * time.Second)
	go Listen(pubsubClient.Channel)
	go pubsubClient.Pub(topicName, "lolwut1")
	go pubsubClient.Pub(topicName, "lolwut2")
	go pubsubClient.Pub(topicName, "lolwut3")
	go pubsubClient.Pub(topicName, "lolwut4")
	go pubsubClient.Pub(topicName, "lolwut5")
	go pubsubClient.Pub(topicName, "lolwut6")

	wg.Wait()
	pubsubClient.Unsub(topicName)
}

func NewPubSubClient(config *Config) *PubSubClient {
	// ipfs daemon --enable-pubsub-experiment
	sh := shell.NewShell("localhost:" + config.Port)
	return &PubSubClient{
		Shell:            sh,
		SubscribedTopics: make(map[string]Topic),
		Channel:          make(chan *shell.Message),
	}
}

func (c *PubSubClient) Pub(topicName string, data string) {
	log.Println("Publishing...")
	if err := c.Shell.PubSubPublish(topicName, data); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished publishing.")
}

func (c *PubSubClient) Sub(topicName string, wg *sync.WaitGroup) {
	subscription, err := c.Shell.PubSubSubscribe(topicName)
	if err != nil {
		log.Fatal(err)
	}
	c.SubscribedTopics[topicName] = Topic{Subscription: subscription}

	for {
		msg, err := subscription.Next()
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
		if err := topic.Subscription.Cancel(); err != nil {
			log.Fatal(err)
		}
		log.Printf("Unsubscribed from the topic: %s", topicName)
	} else {
		log.Printf("There is no subscription for the topic: %s", topicName)
	}
}
