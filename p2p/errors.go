package p2p

import "errors"

var (
	ErrAlreadySubscribed        = errors.New("client has already subscribed to the topic")
	ErrSubscriptionDoesNotExist = errors.New("client is not subscribed to the topic")
	ErrPublishFailed            = errors.New("client could not publish to the subscribed topic")
)
