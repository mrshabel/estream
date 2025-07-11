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
	for idx := range orderCount {
		order := core.Order{ID: int64(idx), Product: getRandomProduct(), CreatedAt: time.Now().UTC()}
		OrderProducer.placeOrder(order)
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

func getRandomProduct() string {
	return products[rand.Intn(len(products))]
}
