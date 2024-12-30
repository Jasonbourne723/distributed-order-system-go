package responses

import (
	"time"
)

type OrderDto struct {
	ID       int64     `json:"id"`
	Amount   int       `json:"amount" `
	Customer string    `json:"customer" `
	Created  time.Time `json:"created" `
	Items    []OrderItem
}

type OrderIdDto struct {
	ID int64 `json:"id"`
}

type OrderItem struct {
	ID       int64  `json:"id"`
	Product  string `json:"product" `
	Quantity int    `json:"quantity" `
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
