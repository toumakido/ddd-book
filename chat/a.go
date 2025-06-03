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
)

type OrderPayment struct {
	Method string
	// メソッドごとの情報
}

type OrderPaymentInformation []struct {
	Method string
	// メソッドごとの情報
}
