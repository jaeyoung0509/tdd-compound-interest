package usecase

import (
	"context"

	"github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/user"
)

// UnitOfWork abstracts transaction management for use cases.
type UnitOfWork interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Payments() payment.Repository
	Users() user.Repository
}
