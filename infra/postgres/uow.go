package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/user"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/repositories"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/sqlc/generated"
	"github.com/jaeyoung0509/compound-interest/usecase"
)

// UnitOfWork manages a single database transaction.
type UnitOfWork struct {
	tx pgx.Tx
}

// NewUnitOfWork creates a new UnitOfWork with an active transaction.
func NewUnitOfWork(ctx context.Context, pool *pgxpool.Pool) (*UnitOfWork, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &UnitOfWork{tx: tx}, nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	return u.tx.Commit(ctx)
}

func (u *UnitOfWork) Rollback(ctx context.Context) error {
	return u.tx.Rollback(ctx)
}

func (u *UnitOfWork) db() generated.DBTX {
	return u.tx
}

func (u *UnitOfWork) Payments() payment.Repository {
	return repositories.NewPaymentRepository(u.db())
}

func (u *UnitOfWork) Users() user.Repository {
	return repositories.NewUserRepository(u.db())
}

// Ensure UnitOfWork implements usecase.UnitOfWork
var _ usecase.UnitOfWork = (*UnitOfWork)(nil)
