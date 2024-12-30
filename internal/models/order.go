package models

import "time"

type Order struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	Amount   int       `json:"amount" gorm:"not null"`
	Customer string    `json:"customer" gorm:"not null"`
	Created  time.Time `json:"created" gorm:"not null"`
	Items    []OrderItem
}
