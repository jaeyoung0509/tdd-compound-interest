package payment

import (
	"context"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

// Repository abstracts persistence for the Payment aggregate.
type Repository interface {
	Get(ctx context.Context, id shared.ID) (*Payment, error)
	Save(ctx context.Context, payment *Payment) error
}
