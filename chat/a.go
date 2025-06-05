package chat

type OrderID string

type Order struct {
	ID          OrderID
	CustomerID  string
	Items       []OrderItem
	Status      OrderStatus
	TotalAmount int64
	Payment     OrderPayment
	Shipping    OrderShipping
	Billing     OrderBilling
}

type OrderItem struct {
	BookID   string
	Quantity int64
}

type OrderStatus string

const (
	OrderStatusCart      OrderStatus = "cart"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderPayment struct {
	Method string
	// その他情報
}

func (o *Order) AllowCancel() bool {
	return o.Status == OrderStatusCart
}
