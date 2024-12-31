package models

import "time"

type Inventory struct {
	ID      int64     `json:"id" gorm:"primary_key"`
	Product string    `json:"product" gorm:"unique not null"`
	Created time.Time `json:"created" gorm:"not null"`
	Stock   int64     `json:"stock" gorm:"not null"`
	Stock1  int64     `json:"stock1" gorm:"not null"`
	Stock2  int64     `json:"stock2" gorm:"not null"`
	Stock3  int64     `json:"stock3" gorm:"not null"`
	Stock4  int64     `json:"stock4" gorm:"not null"`
	Stock5  int64     `json:"stock5" gorm:"not null"`
	Stock6  int64     `json:"stock6" gorm:"not null"`
	Stock7  int64     `json:"stock7" gorm:"not null"`
	Stock8  int64     `json:"stock8" gorm:"not null"`
	Stock9  int64     `json:"stock9" gorm:"not null"`
	Stock10 int64     `json:"stock10" gorm:"not null"`
	Stock11 int64     `json:"stock11" gorm:"not null"`
	Stock12 int64     `json:"stock12" gorm:"not null"`
	Stock13 int64     `json:"stock13" gorm:"not null"`
	Stock14 int64     `json:"stock14" gorm:"not null"`
	Stock15 int64     `json:"stock15" gorm:"not null"`
	Stock16 int64     `json:"stock16" gorm:"not null"`
	Stock17 int64     `json:"stock17" gorm:"not null"`
	Stock18 int64     `json:"stock18" gorm:"not null"`
	Stock19 int64     `json:"stock19" gorm:"not null"`
	Stock20 int64     `json:"stock20" gorm:"not null"`
}
