// this package contains code that consumes order and inventory updates from the kafka broker and fire a notification
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mrshabel/estream/core"
)

func sendNotification(order core.Order) error {
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("New order fulfilled: %s\n", order)
	return nil
}

func main() {
	consumer, err := NewNotificationConsumer()
	if err != nil {
		log.Fatalf("failed to create inventory consumer: %v\n", err)
	}
	defer consumer.Close()

	// start consumer with callback
	if err := consumer.Start(sendNotification); err != nil {
		log.Fatalf("inventory consumer failed: %v\n", err)
	}
}
