package core

import (
	"fmt"
	"time"
)

type Order struct {
	ID        int64     `json:"id"`
	Product   string    `json:"product"`
	CreatedAt time.Time `json:"createdAt"`
}

func (o *Order) String() string {
	return fmt.Sprintf("Order %d: %v", o.ID, o.Product)
}
