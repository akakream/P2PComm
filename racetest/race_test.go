package racetest

import (
	"context"
	"log"
	"testing"

	p2p "github.com/akakream/sailorsailor/p2p"
)

func TestRaceConditionLibp2p(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := p2p.NewLibP2PClient(ctx, false)

	topicName1 := "randomTopic1"
	topicName2 := "randomTopic2"

	go func() {
		err := c.Sub(topicName1)
		if err != nil {
			log.Printf("Sub failed with error %v", err)
		}
	}()
	go func() {
		err := c.Sub(topicName2)
		if err != nil {
			log.Printf("Sub failed with error %v", err)
		}
	}()

	c.Shutdown()
}
