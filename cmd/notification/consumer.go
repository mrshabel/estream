package main

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mrshabel/estream/core"
)

// topics to consume from
var (
	NotificationConsumerTopics  = []string{"order.placed", "inventory.reserved", "inventory.insufficient"}
	NotificationConsumerGroupID = "notification-service-group"
)

type NotificationConsumer struct {
	Consumer *kafka.Consumer
}

type Callback func(order core.Order) error

func NewNotificationConsumer() (*NotificationConsumer, error) {
	consumer, err := core.NewConsumer(NotificationConsumerTopics, NotificationConsumerGroupID)
	if err != nil {
		return nil, err
	}

	// subscribe to topics
	if err = consumer.SubscribeTopics(NotificationConsumerTopics, nil); err != nil {
		return nil, err
	}

	return &NotificationConsumer{
		Consumer: consumer,
	}, nil
}

// Start consumes messages from the broker in a blocking manner. A callback should be passed to the consumer to process the received message
func (c *NotificationConsumer) Start(callback Callback) error {
	log.Println("Notification consumer waiting for messages...")
	for {
		// poll in 100 milliseconds
		ev := c.Consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			// parse message
			var order core.Order
			if err := json.Unmarshal(e.Value, &order); err != nil {
				log.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			// process on callback
			if err := callback(order); err != nil {
				log.Printf("Consumer callback processing error: %v\n", err)
				continue
			}

			// commit
			if _, err := c.Consumer.CommitMessage(e); err != nil {
				log.Printf("Failed to commit message (%s): %v\n", e, err)
			}
		case kafka.Error:
			if e.IsFatal() {
				return e
			}
		// no event received
		case nil:
			continue
		}
	}
}

func (c *NotificationConsumer) Close() error {
	return c.Consumer.Close()
}
