package test

import (
	"testing"

	p2p "github.com/akakream/sailorsailor/p2p"
)

func TestRaceCondition(t *testing.T) {
	c := p2p.NewLibP2PClient()

	topicName1 := "randomTopic1"
	topicName2 := "randomTopic2"

	go c.Sub(topicName1)
	go c.Sub(topicName2)

	c.Shutdown()
}
