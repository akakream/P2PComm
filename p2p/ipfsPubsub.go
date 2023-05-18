package p2p

import (
	"fmt"
	"log"
	"time"

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

func Run() {
	sh := shell.NewShell("localhost:" + "5001")
	_ = sh

	topicName := "babuska1"
	pubsubClient := NewIpfsP2PClient(&Config{Port: "5001"})
	go pubsubClient.Sub(topicName)
	time.Sleep(5 * time.Second)
	go pubsubClient.listen(pubsubClient.Channel)
	go pubsubClient.Pub(topicName, "lolwut1")
	go pubsubClient.Pub(topicName, "lolwut2")
	go pubsubClient.Pub(topicName, "lolwut3")
	go pubsubClient.Pub(topicName, "lolwut4")
	go pubsubClient.Pub(topicName, "lolwut5")
	go pubsubClient.Pub(topicName, "lolwut6")

	pubsubClient.Unsub(topicName)
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

func (c *IpfsP2PClient) Pub(topicName string, data string) {
	log.Println("Publishing...")
	if err := c.Shell.PubSubPublish(topicName, data); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished publishing.")
}

func (c *IpfsP2PClient) Sub(topicName string) {
	subscription, err := c.Shell.PubSubSubscribe(topicName)
	if err != nil {
		log.Fatal(err)
	}
	c.SubscribedTopics[topicName] = IpfsP2PTopic{Subscription: subscription}

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
