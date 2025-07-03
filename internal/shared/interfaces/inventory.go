package interfaces

import (
	"context"

	"github.com/toumakido/ddd-book/internal/shared/domain/model"
)

type InventoryService interface {
	CheckAvailability(ctx context.Context, bookID model.BookID, quantity int64) (bool, error)
}
