package mq

import (
	"errors"
	"math/rand"
)

var (
	errMQPublic = errors.New("error public message")
	errMQSub    = errors.New("error subscribe message")
)

type MessageQueuePubliser interface {
	Public(topic string, message interface{}) error
}

type MessageQueueSubscriber interface {
	Subscribe(topic string, handler func(msg interface{}) error) error
}

type MessageQueue interface {
	MessageQueuePubliser
	MessageQueueSubscriber
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

func NewMockMQ() MessageQueue {
	return &mockMessageQueue{
		rand: rand.New(rand.NewSource(0)),
	}
}

type mockSubscriber struct {
}

func (s *mockSubscriber) Subscribe(topic string, handler func(msg interface{}) error) error {
	handler("Mocking message")
	return nil
}

func NewMockSubscriber() MessageQueueSubscriber {
	return &mockSubscriber{}
}

type mockPublisher struct {
}

func (p *mockPublisher) Public(topic string, message interface{}) error {
	return nil
}

func NewMockPublisher() MessageQueuePubliser {
	return &mockPublisher{}
}
