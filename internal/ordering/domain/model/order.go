package model

import (
	"time"

	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
)

type OrderID string

type Order struct {
	ID          OrderID
	CustomerID  sharedmodel.CustomerID
	Items       []OrderItem
	Status      OrderStatus
	TotalAmount int64
	Shipping    OrderShipping
	Payment     OrderPayment
	Billing     OrderBilling
}

type OrderItem struct {
	BookID   sharedmodel.BookID
	Quantity int
	Amount   int64
}

type OrderStatus string

const (
	OrderStatusCart      OrderStatus = "cart"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCanceled  OrderStatus = "canceled"
)

type OrderShipping struct {
	Address     sharedmodel.Address
	Method      string
	TrackingID  string
	ShippedAt   *time.Time
	DeliveredAt *time.Time
}
