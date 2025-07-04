package model

import (
	"errors"
	"time"

	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
)

var (
	ErrInvalidQuantity       = errors.New("quantity must be positive")
	ErrOrderAlreadyConfirmed = errors.New("order has already been confirmed")
)

func (o *Order) AddItem(bookID sharedmodel.BookID, quantity int64, amount int64) error {
	if o.status != OrderStatusCart {
		return ErrOrderAlreadyConfirmed
	}

	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	// 既存アイテムの数量を更新
	for i, item := range o.items {
		if item.BookID == bookID {
			o.items[i].Quantity += quantity
			o.recalculateTotalAmount()
			o.updatedAt = time.Now()
			return nil
		}
	}

	// 新しいアイテムを追加
	o.items = append(o.items, OrderItem{
		BookID:   bookID,
		Quantity: quantity,
		Amount:   amount,
	})

	o.recalculateTotalAmount()
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) recalculateTotalAmount() {
	var total int64
	for _, item := range o.items {
		total += item.Amount * int64(item.Quantity)
	}
	o.totalAmount = total
}

func (o *Order) Items() []OrderItem {
	return o.items
}

func (o *Order) Confirm() error {
	if o.status != OrderStatusCart {
		return ErrOrderAlreadyConfirmed
	}

	o.status = OrderStatusConfirmed
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) SetShipping(address sharedmodel.Address, fee int64, method string) {
	o.shipping = OrderShipping{
		Address: address,
		Fee:     fee,
		Method:  method,
	}
}

func (o *Order) PaymentAmount() int64 {
	return o.totalAmount + o.shipping.Fee
}

func (o *Order) SetPayment(method string) {
	o.payment = OrderPayment{
		Method: method,
		Amount: o.PaymentAmount(),
	}
	o.updatedAt = time.Now()
}
