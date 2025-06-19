package repository

import (
	"context"

	"github.com/toumakido/ddd-book/internal/ordering/domain/model"
	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
)

type OrderRepository interface {
	GetOrderByCustomerIDForUpdate(ctx context.Context, customerID sharedmodel.CustomerID) (*model.Order, error)
	Save(ctx context.Context, order *model.Order) error
}
