package model

import (
	"time"

	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
)

// 集約ルート
type Order struct {
	id          OrderID
	customerID  sharedmodel.CustomerID
	items       []OrderItem
	status      OrderStatus
	totalAmount int64
	shipping    OrderShipping
	payment     OrderPayment
	billing     OrderBilling
	createdAt   time.Time
	updatedAt   time.Time
}

// 値オブジェクト
type OrderID string

type OrderItem struct {
	BookID   sharedmodel.BookID
	Quantity int64
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
	Fee         int64
	Method      string
	TrackingID  *string
	ShippedAt   *time.Time
	DeliveredAt *time.Time
}

type OrderPayment struct {
	Method      string
	Amount      int64
	ProcessedAt *time.Time
}

type OrderBilling struct {
	Address     sharedmodel.Address
	BillingID   string
	TotalAmount int64
}

func NewOrder(customerID sharedmodel.CustomerID) *Order {
	return &Order{
		id:         NewOrderID(),
		customerID: customerID,
		items:      make([]OrderItem, 0),
		status:     OrderStatusCart,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}
}

func NewOrderID() OrderID {
	return OrderID(sharedmodel.NewID())
}
