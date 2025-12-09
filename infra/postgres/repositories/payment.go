package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Get(ctx context.Context, id shared.ID) (*payment.Payment, error) {
	return nil, nil
}

func (r *PaymentRepository) Save(ctx context.Context, p *payment.Payment) error {
	return nil
}
