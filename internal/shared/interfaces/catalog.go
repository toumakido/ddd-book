package interfaces

import (
	"context"

	"github.com/toumakido/ddd-book/internal/shared/domain/model"
)

type CatalogService interface {
	GetBookPrice(ctx context.Context, bookID model.BookID) (int64, error)
}
