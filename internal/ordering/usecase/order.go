package usecase

import (
	"context"

	"github.com/toumakido/ddd-book/internal/ordering/domain/model"
	"github.com/toumakido/ddd-book/internal/ordering/domain/repository"
	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
	sharedrepository "github.com/toumakido/ddd-book/internal/shared/domain/repository"
)

type orderUsecase struct {
	orderRepository repository.OrderRepository
}

type AddItemToCartParams struct {
	CustomerID sharedmodel.CustomerID
	BookID     string
	Quantity   int
}

func (u *orderUsecase) AddItemToCart(ctx context.Context, params AddItemToCartParams) error {
	// TODO: トランザクション
	order, err := u.orderRepository.GetOrderByCustomerIDForUpdate(ctx, params.CustomerID)
	if err != nil {
		if err == sharedrepository.ErrNotFound {
			order = model.NewOrder(params.CustomerID)
		} else {
			return err
		}
	}
	order.AddItem(sharedmodel.BookID(params.BookID), int64(params.Quantity), 0) // TODO: 金額は外部から取得する
	if err := u.orderRepository.Save(ctx, order); err != nil {
		return err
	}
	return nil
}
