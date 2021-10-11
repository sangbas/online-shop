package entity

import (
	"time"
)

type Order struct {
	ID            string     `db:"id"`
	UserID        string     `db:"user_id"`
	AddressID     string     `db:"address_id"`
	OrderDate     *time.Time `db:"order_date"`
	PaymentDate   *time.Time `db:"payment_date"`
	DeliveredDate *time.Time `db:"delivered_date"`
	Status        string     `db:"status"`
	Amount        float64    `db:"amount"`
	OrderDetails  []OrderDetail
}

type OrderDetail struct {
	ID        string  `db:"id"`
	OrderID   string  `db:"order_id"`
	ProductID int64   `db:"product_id"`
	Price     float64 `db:"price"`
	Quantity  int32   `db:"quantity"`
}

type CompleteOrder struct {
	ID            string     `db:"id"`
	UserID        string     `db:"user_id"`
	AddressID     string     `db:"address_id"`
	OrderDate     *time.Time `db:"order_date"`
	PaymentDate   *time.Time `db:"payment_date"`
	DeliveredDate *time.Time `db:"delivered_date"`
	Status        string     `db:"status"`
	Amount        float64    `db:"amount"`
	ProductName   string     `db:"name"`
	Price         float64    `db:"price"`
	Quantity      int32      `db:"quantity"`
}
