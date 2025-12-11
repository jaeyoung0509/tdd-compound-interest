package testhelper

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/domain/user"
)

// PaymentOption mutates a Payment before persisting it.
type PaymentOption func(*payment.Payment) error

// WithPaidAt marks a payment as paid before insertion.
func WithPaidAt(paidAt time.Time) PaymentOption {
	return func(p *payment.Payment) error {
		return p.Pay(paidAt)
	}
}

// CreateTestUser inserts a user row with deterministic timestamps when provided.
func CreateTestUser(ctx context.Context, pool *pgxpool.Pool, name string, now time.Time) (user.ID, error) {
	if now.IsZero() {
		now = time.Now()
	}

	u, err := user.New(name, now)
	if err != nil {
		return user.ID{}, err
	}

	_, err = pool.Exec(ctx,
		"INSERT INTO users (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		u.ID().Value().String(), u.Name(), u.CreatedAt(), u.UpdatedAt(),
	)
	return u.ID(), err
}

// CreateTestPayment builds a domain payment, applies options, and persists it.
func CreateTestPayment(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID user.ID,
	amount money.Money,
	dueDate time.Time,
	now time.Time,
	opts ...PaymentOption,
) (shared.ID, error) {
	if now.IsZero() {
		now = time.Now()
	}

	p, err := payment.New(userID, amount, dueDate, now)
	if err != nil {
		return shared.ID{}, err
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return shared.ID{}, err
		}
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO payments (id, user_id, amount, currency, due_date, paid_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		p.ID().String(), p.UserID().Value().String(),
		p.Amount().Amount(), string(p.Amount().Currency()),
		p.DueDate(), p.PaidAt(), string(p.Status()),
		p.CreatedAt(), p.UpdatedAt(),
	)

	return p.ID(), err
}
