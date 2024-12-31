package requests

type CreateOrderDto struct {
	OrderId  int64  `json:"order_id" `
	Amount   int    `json:"amount" `
	Customer string `json:"customer" `
	Items    []OrderItem
}

type OrderItem struct {
	Product  string `json:"product" `
	Quantity int    `json:"quantity" `
	Price    int    `json:"price"`
	Amount   int    `json:"amount"`
}
