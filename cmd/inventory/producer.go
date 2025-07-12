package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mrshabel/estream/core"
)

// constants
const (
	clientID             = "inventory-service"
	successProducerTopic = "inventory.reserved"
	failedProducerTopic  = "inventory.insufficient"
)

// inventory statuses
type InventoryStatus int

const (
	Sufficient InventoryStatus = iota
	Insufficient
)

type InventoryProducer struct {
	Producer          *kafka.Producer
	ReservedTopic     string
	InsufficientTopic string
	// error delivery channel
	DeliveryChan chan kafka.Event
}

func NewInventoryProducer() (*InventoryProducer, error) {
	producer, err := core.NewProducer(clientID)
	if err != nil {
		return nil, err
	}

	return &InventoryProducer{
		Producer:          producer,
		ReservedTopic:     successProducerTopic,
		InsufficientTopic: failedProducerTopic,
		DeliveryChan:      make(chan kafka.Event, 1000),
	}, nil
}

// processInventory simulates inventory checks and based on the status of the inventory: sufficient or insufficient
func (i *InventoryProducer) processInventory(order core.Order, status InventoryStatus) error {
	// encode order payload into bytes
	payload, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order payload: %w", err)
	}
	var data *kafka.Message
	switch status {
	case Sufficient:
		data = &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &i.ReservedTopic, Partition: kafka.PartitionAny},
			Value:          payload,
		}
	case Insufficient:
		data = &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &i.InsufficientTopic, Partition: kafka.PartitionAny},
			Value:          payload,
		}
	default:
		return fmt.Errorf("invalid inventory status received: %v", status)
	}

	if err := i.Producer.Produce(data, i.DeliveryChan); err != nil {
		return err
	}

	// wait for delivery report
	e := <-i.DeliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return fmt.Errorf("failed to reserve: %w", msg.TopicPartition.Error)
	}

	return nil
}
