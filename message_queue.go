package main

import (
	"errors"
	"math/rand"
)

var (
	errMQPublic = errors.New("error public message")
	errMQSub    = errors.New("error subscribe message")
)

type MessageQueue interface {
	Public(topic string, message interface{}) error
	Subscribe(topic string, handler func(msg interface{}) error) error
}

type mockMessageQueue struct {
	rand *rand.Rand
}

func (mq *mockMessageQueue) Public(topic string, message interface{}) error {
	r := mq.rand.Intn(100)
	if r > 10 {
		return nil
	}
	return errMQPublic
}

func (mq *mockMessageQueue) Subscribe(topic string, handler func(msg interface{}) error) error {
	r := mq.rand.Intn(100)
	if r > 10 {
		handler("Mocking message")
		return nil
	}
	return errMQSub
}

func newMockMQ() MessageQueue {
	return &mockMessageQueue{
		rand: rand.New(rand.NewSource(0)),
	}
}
