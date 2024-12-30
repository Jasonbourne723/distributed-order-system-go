package requests

import "time"

type CreateOrderDto struct {
	OrderId  int64     `json:"order_id" `
	Amount   int       `json:"amount" `
	Customer string    `json:"customer" `
	Created  time.Time `json:"created" `
	Items    []OrderItem
}

type OrderItem struct {
	Product  string `json:"product" `
	Quantity int    `json:"quantity" `
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
