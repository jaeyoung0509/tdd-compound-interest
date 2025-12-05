package payment

import (
	"context"

	dp "github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

// Service orchestrates payment accrual with injected dependencies.
type Service struct {
	repo         dp.Repository
	clock        dp.Clock
	rateProvider dp.DailyRateProvider
}

func NewService(repo dp.Repository, clock dp.Clock, rateProvider dp.DailyRateProvider) *Service {
	return &Service{
		repo:         repo,
		clock:        clock,
		rateProvider: rateProvider,
	}
}

// AccruePayment loads a payment, accrues interest using the injected collaborators, and persists the result.
func (s *Service) AccruePayment(ctx context.Context, id shared.ID) (*dp.Payment, error) {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.AccrueInterestWith(s.clock, s.rateProvider); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
