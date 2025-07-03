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
	// 他のメソッドも追加予定
}
type orderUseCaseImpl struct {
	orderRepository  repository.OrderRepository
	catalogService   interfaces.CatalogService
	inventoryService interfaces.InventoryService
}

// コンストラクタ
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	catalogSvc interfaces.CatalogService,
	inventorySvc interfaces.InventoryService,
) OrderUseCase {
	return &orderUseCaseImpl{
		orderRepository:  orderRepo,
		catalogService:   catalogSvc,
		inventoryService: inventorySvc,
	}
}

type AddItemToCartParams struct {
	CustomerID sharedmodel.CustomerID
	BookID     string
	Quantity   int64
}

func (u *orderUseCaseImpl) AddItemToCart(ctx context.Context, params AddItemToCartParams) (*model.Order, error) {
	bookID := sharedmodel.BookID(params.BookID)

	tx, err := u.orderRepository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	order, err := u.orderRepository.GetOrderByCustomerIDForUpdate(ctx, tx, params.CustomerID)
	if err != nil {
		if err == sharedrepository.ErrNotFound {
			order = model.NewOrder(params.CustomerID)
		} else {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
	}

	// 在庫確認（在庫サービス経由）
	available, err := u.inventoryService.CheckAvailability(ctx, bookID, params.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to check inventory: %w", err)
	}
	if !available {
		return nil, errors.New("item is out of stock or not enough quantity available")
	}

	// カタログサービスから価格を取得
	bookPrice, err := u.catalogService.GetBookPrice(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book price: %w", err)
	}

	if err := order.AddItem(bookID, params.Quantity, bookPrice); err != nil {
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
