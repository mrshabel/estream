package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	Topic = "orders.placed"
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "orderConsumer",
		"auto.offset.reset": "smallest"},
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}
	defer consumer.Close()

	if err = consumer.Subscribe(Topic, nil); err != nil {
		log.Fatalf("Failed to subscribe to topic(s): %s\n", err)
	}

	log.Println("Waiting for messages...")
	for {
		// poll in 100 milliseconds
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Consumed message from topic %s: %s\n", Topic, string(e.Value))
			// commit message offset
			consumer.Commit()
		case kafka.Error:
			fmt.Printf("consumer kafka error%+v\n", e)
			if e.IsFatal() {
				return
			}
		// no event received
		case nil:
			continue
		}
	}
}
