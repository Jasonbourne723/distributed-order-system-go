package models

// OrderItem represents an item in an order
type OrderItem struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	OrderId  int64  `json:"order_id" gorm:"index"`
	Product  string `json:"product" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Price    int    `json:"price" gorm:"not null default 0"`
	Amount   int    `json:"price" gorm:"not null default 0"`
}
