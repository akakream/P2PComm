package main

import (

	// "context"

	"fmt"
	"sync"
	"time"

	pubsub "github.com/akakream/sailorsailor/pubsub"
	shell "github.com/ipfs/go-ipfs-api"
)

func listen(channel chan *shell.Message) {
	for {
		select {
		case msg := <-channel:
			fmt.Println(string(msg.Data))
			fmt.Println(msg.TopicIDs)
		default:
		}
	}
}

func main() {
	sh := shell.NewShell("localhost:" + "5001")
	_ = sh

	topicName := "babuska1"
	pubsubClient := pubsub.NewPubSubClient(&pubsub.Config{Port: "5001"})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go pubsubClient.Sub(topicName, wg)
	time.Sleep(5 * time.Second)
	go listen(pubsubClient.Channel)
	go pubsubClient.Pub(topicName, "lolwut")
	go pubsubClient.Pub(topicName, "lolwut2")
	go pubsubClient.Pub(topicName, "lolwut3")
	go pubsubClient.Pub(topicName, "lolwut4")
	go pubsubClient.Pub(topicName, "lolwut5")
	go pubsubClient.Pub(topicName, "lolwut6")

	wg.Wait()
	pubsubClient.Unsub(topicName)
}
