package payment

import (
	"context"
	"sync"

	dp "github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

// InMemoryPaymentRepo is a simple fake repository for tests and local usage.
type InMemoryPaymentRepo struct {
	mu       sync.Mutex
	store    map[shared.ID]*dp.Payment
	GetErr   error
	SaveErr  error
	saveHits int
}

func NewInMemoryPaymentRepo() *InMemoryPaymentRepo {
	return &InMemoryPaymentRepo{
		store: make(map[shared.ID]*dp.Payment),
	}
}

// Seed primes the repository with given payments.
func (r *InMemoryPaymentRepo) Seed(payments ...*dp.Payment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range payments {
		r.store[p.ID()] = p
	}
}

func (r *InMemoryPaymentRepo) Get(ctx context.Context, id shared.ID) (*dp.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.GetErr != nil {
		return nil, r.GetErr
	}
	p, ok := r.store[id]
	if !ok {
		return nil, dp.ErrPaymentNotFound
	}
	return p, nil
}

func (r *InMemoryPaymentRepo) Save(ctx context.Context, payment *dp.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.SaveErr != nil {
		return r.SaveErr
	}
	r.store[payment.ID()] = payment
	r.saveHits++
	return nil
}

func (r *InMemoryPaymentRepo) SaveCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.saveHits
}
