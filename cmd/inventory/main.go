// this package contains code that consumes placed orders from the kafka broker and process their inventory statuses
package main

import (
	"log"
)

func main() {
	// create order placer
	producer, err := NewInventoryProducer()
	if err != nil {
		log.Fatalf("failed to create inventory producer: %v\n", err)
	}
	// close delivery channel when done
	defer close(producer.DeliveryChan)

	consumer, err := NewInventoryConsumer()
	if err != nil {
		log.Fatalf("failed to create inventory consumer: %v\n", err)
	}
	defer consumer.Close()

	// start consumer with callback
	if err := consumer.Start(producer.processInventory); err != nil {
		log.Fatalf("inventory consumer failed: %v\n", err)
	}

}
