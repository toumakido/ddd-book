package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Tx interface {
	Commit() error
	Rollback() error
	RollbackUnlessCommitted()
}

type Transactionable interface {
	BeginTx(ctx context.Context) (Tx, error)
}
