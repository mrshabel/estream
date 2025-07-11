package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mrshabel/estream/core"
)

// constants
const (
	clientID      = "order-service"
	producerTopic = "order.placed"
)

type OrderProducer struct {
	Producer *kafka.Producer
	Topic    string
	// error delivery channel
	DeliveryChan chan kafka.Event
}

func NewOrderProducer() (*OrderProducer, error) {
	producer, err := core.NewProducer(clientID)
	if err != nil {
		return nil, err
	}

	return &OrderProducer{
		Producer:     producer,
		Topic:        producerTopic,
		DeliveryChan: make(chan kafka.Event, 1000),
	}, nil
}

// placeOrder places a new order by sending the order to the kafaka broker
func (o *OrderProducer) placeOrder(order core.Order) error {
	// encode order payload into bytes
	payload, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order payload: %w", err)
	}

	if err := o.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &o.Topic, Partition: kafka.PartitionAny},
			Value:          payload,
		},
		o.DeliveryChan,
	); err != nil {
		return err
	}

	// wait for delivery report
	e := <-o.DeliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return fmt.Errorf("failed to place order: %w", msg.TopicPartition.Error)
	}

	return nil
}
