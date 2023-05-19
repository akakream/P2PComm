package p2p

import "errors"

var (
	ErrAlreadySubscribed = errors.New("client has already subscribed to the topic")
)
