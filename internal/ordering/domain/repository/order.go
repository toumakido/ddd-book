package repository

import (
	"context"

	"github.com/toumakido/ddd-book/internal/ordering/domain/model"
	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
	sharedrepository "github.com/toumakido/ddd-book/internal/shared/domain/repository"
)

type OrderRepository interface {
	sharedrepository.Transactionable
	GetOrderByCustomerIDForUpdate(ctx context.Context, tx sharedrepository.Tx, customerID sharedmodel.CustomerID) (*model.Order, error)
	Save(ctx context.Context, tx sharedrepository.Tx, order *model.Order) error
}
