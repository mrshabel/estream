package core

import (
	"fmt"
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Product   string    `json:"product"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

func (o Order) String() string {
	return fmt.Sprintf("Order %d: %v", o.ID, o.Product)
}
