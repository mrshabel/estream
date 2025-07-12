// this package contains code for generating random orders and sending it to the broker for downstream consumers to continue processing
package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/mrshabel/estream/core"
)

func main() {
	// cli flags
	var orderCount int

	flag.IntVar(&orderCount, "count", 100, "Number of random orders to place")
	flag.Parse()

	// create order placer
	OrderProducer, err := NewOrderProducer()
	if err != nil {
		log.Fatalf("failed to create order producer: %v\n", err)
	}
	// close delivery channel when done
	defer close(OrderProducer.DeliveryChan)

	// produce bulk messages
	log.Println("Starting bulk order placing...")
	for idx := range orderCount {
		order := getRandomOrder(idx)
		OrderProducer.placeOrder(order)
		log.Printf("Order placed successfully: %s\n", order)
		// throttle requests
		time.Sleep(2 * time.Second)
	}

}

var products = []string{
	"Virgil",
	"Adidas Ultraboost",
	"Puma",
	"New Balance",
	"Reebok Classic",
	"FTY",
}

func getRandomOrder(idx int) core.Order {
	product := products[rand.Intn(len(products))]
	return core.Order{ID: idx, Product: product, Quantity: uint(rand.Intn(100)), CreatedAt: time.Now().UTC()}
}
