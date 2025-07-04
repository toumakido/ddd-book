package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/toumakido/ddd-book/internal/ordering/domain/model"
	"github.com/toumakido/ddd-book/internal/ordering/domain/repository"
	sharedmodel "github.com/toumakido/ddd-book/internal/shared/domain/model"
	sharedrepository "github.com/toumakido/ddd-book/internal/shared/domain/repository"
	"github.com/toumakido/ddd-book/internal/shared/interfaces"
)

type OrderUseCase interface {
	AddItemToCart(ctx context.Context, params AddItemToCartParams) (*model.Order, error)
}
type orderUseCaseImpl struct {
	orderRepository  repository.OrderRepository
	paymentService   repository.PaymentService
	catalogService   interfaces.CatalogService
	inventoryService interfaces.InventoryService
}

// コンストラクタ
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentService,
	catalogSvc interfaces.CatalogService,
	inventorySvc interfaces.InventoryService,
) OrderUseCase {
	return &orderUseCaseImpl{
		orderRepository:  orderRepo,
		paymentService:   paymentRepo,
		catalogService:   catalogSvc,
		inventoryService: inventorySvc,
	}
}

type AddItemToCartParams struct {
	CustomerID sharedmodel.CustomerID
	BookID     sharedmodel.BookID
	Quantity   int64
}

func (u *orderUseCaseImpl) AddItemToCart(ctx context.Context, params AddItemToCartParams) (*model.Order, error) {
	tx, err := u.orderRepository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	order, err := u.orderRepository.GetByCustomerIDForUpdate(ctx, tx, params.CustomerID)
	if err != nil {
		if err == sharedrepository.ErrNotFound {
			order = model.NewOrder(params.CustomerID)
		} else {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
	}

	available, err := u.inventoryService.CheckAvailability(ctx, params.BookID, params.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to check inventory: %w", err)
	}
	if !available {
		return nil, errors.New("item is out of stock or not enough quantity available")
	}

	bookPrice, err := u.catalogService.GetBookPrice(ctx, params.BookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book price: %w", err)
	}

	if err := order.AddItem(params.BookID, params.Quantity, bookPrice); err != nil {
		return nil, fmt.Errorf("failed to add item to order: %w", err)
	}
	if err := u.orderRepository.Save(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return order, nil
}

type ConfirmOrderParams struct {
	OrderID          model.OrderID
	ShipppingAddress sharedmodel.Address
	ShippingFee      int64
	ShippingMethod   string
	PaymentMethod    string
}

func (u *orderUseCaseImpl) ConfirmOrder(ctx context.Context, params ConfirmOrderParams) (*model.Order, error) {
	tx, err := u.orderRepository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	order, err := u.orderRepository.GetByIDForUpdate(ctx, tx, params.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	reservationIDs := make([]string, 0, len(order.Items()))
	for _, item := range order.Items() {
		available, err := u.inventoryService.CheckAvailability(ctx, item.BookID, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to check inventory: %w", err)
		}
		if !available {
			return nil, errors.New("item is out of stock or not enough quantity available")
		}
		reservationID, err := u.inventoryService.Reserve(ctx, item.BookID, item.Quantity)
		if err != nil {
			// 既に予約したものをキャンセル
			for _, id := range reservationIDs {
				if err := u.inventoryService.CancelReservation(ctx, id); err != nil {
					fmt.Printf("failed to cancel reservation: %s, id: %s", err, id)
				}
			}
			return nil, fmt.Errorf("failed to reserve inventory: %w", err)
		}
		reservationIDs = append(reservationIDs, reservationID)
	}

	if err := order.Confirm(); err != nil {
		return nil, fmt.Errorf("failed to confirm order: %w", err)
	}

	order.SetShipping(params.ShipppingAddress, params.ShippingFee, params.ShippingMethod)

	// 決済問い合わせ
	if err := u.paymentService.RequestPayment(ctx, params.PaymentMethod, order.PaymentAmount()); err != nil {
		// 支払い失敗時に在庫予約をキャンセル
		for _, id := range reservationIDs {
			if err := u.inventoryService.CancelReservation(ctx, id); err != nil {
				fmt.Printf("failed to cancel reservation: %s, id: %s", err, id)
			}
		}
		return nil, fmt.Errorf("failed to request payment: %w", err)
	}

	order.SetPayment(params.PaymentMethod)

	if err := u.orderRepository.Save(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}
