package repository

import (
	"context"
)

type PaymentService interface {
	RequestPayment(ctx context.Context, method string, amount int64) error
}
