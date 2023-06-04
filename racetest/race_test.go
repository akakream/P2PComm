package racetest

import (
	"context"
	"testing"

	p2p "github.com/akakream/sailorsailor/p2p"
)

func TestRaceConditionLibp2p(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := p2p.NewLibP2PClient(ctx)

	topicName1 := "randomTopic1"
	topicName2 := "randomTopic2"

	go c.Sub(topicName1)
	go c.Sub(topicName2)

	c.Shutdown()
}

func TestRaceConditionIPFS(t *testing.T) {
	c := p2p.NewIpfsP2PClient(&p2p.Config{Port: "5001"})

	topicName1 := "randomTopic1"
	topicName2 := "randomTopic2"

	go c.Sub(topicName1)
	go c.Sub(topicName2)

	c.Shutdown()
}
