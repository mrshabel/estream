package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mrshabel/estream/core"
)

// topics to consume from
var (
	InventoryConsumerTopics  = []string{"order.placed"}
	InventoryConsumerGroupID = "inventory-service-group"
)

type InventoryConsumer struct {
	Consumer *kafka.Consumer
	// OrderMessagesCh chan *core.Order
}

// consumer callback function
type Callback func(core.Order, InventoryStatus) error

func NewInventoryConsumer() (*InventoryConsumer, error) {
	consumer, err := core.NewConsumer(InventoryConsumerTopics, InventoryConsumerGroupID)
	if err != nil {
		return nil, err
	}

	// subscribe to topics
	if err = consumer.SubscribeTopics(InventoryConsumerTopics, nil); err != nil {
		return nil, err
	}

	return &InventoryConsumer{
		Consumer: consumer,
		// buffered channel with size 1000
		// OrderMessagesCh: make(chan *core.Order, 1000),
	}, nil
}

// Start consumes messages from the broker in a blocking manner. A callback should be passed to the consumer to process the received message
func (c *InventoryConsumer) Start(callback Callback) error {
	log.Println("Inventory consumer waiting for messages...")
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
			status := simulateInventoryUpdate()
			if err := callback(order, status); err != nil {
				log.Printf("Consumer callback processing error: %v\n", err)
				continue
			}
			log.Printf("Inventory processed order: (%s) with status: %v\n", order, status)

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

func (c *InventoryConsumer) Close() error {
	return c.Consumer.Close()
}

// helper methods
func simulateInventoryUpdate() InventoryStatus {
	delay := time.Duration(1+rand.Intn(4)) * time.Second
	time.Sleep(delay)
	// insufficient if delay is longer
	if delay > 2 {
		return Insufficient
	}
	return Sufficient
}
