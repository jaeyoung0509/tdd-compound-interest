package payment

import (
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

const (
	EventPaymentOverdueAccrued = "payment.overdue_accrued"
	EventPaymentPaid           = "payment.paid"
)

type OverdueAccrued struct {
	PaymentID       string    `json:"payment_id"`
	UserID          string    `json:"user_id"`
	DaysOverdue     int       `json:"days_overdue"`
	PenaltyAmount   string    `json:"penalty_amount"`
	PenaltyCurrency string    `json:"penalty_currency"`
	CalculatedAt    time.Time `json:"calculated_at"`
	OccurredAtTime  time.Time `json:"occurred_at"`
}

func (e OverdueAccrued) EventType() string {
	return EventPaymentOverdueAccrued
}

func (e OverdueAccrued) AggregateType() string {
	return "payment"
}

func (e OverdueAccrued) AggregateID() string {
	return e.PaymentID
}

func (e OverdueAccrued) OccurredAt() time.Time {
	return e.OccurredAtTime
}

var _ shared.DomainEvent = OverdueAccrued{}

type PaymentPaid struct {
	PaymentID      string    `json:"payment_id"`
	UserID         string    `json:"user_id"`
	PaidAt         time.Time `json:"paid_at"`
	OccurredAtTime time.Time `json:"occurred_at"`
}

func (e PaymentPaid) EventType() string {
	return EventPaymentPaid
}

func (e PaymentPaid) AggregateType() string {
	return "payment"
}

func (e PaymentPaid) AggregateID() string {
	return e.PaymentID
}

func (e PaymentPaid) OccurredAt() time.Time {
	return e.OccurredAtTime
}

var _ shared.DomainEvent = PaymentPaid{}

func newOverdueAccruedEvent(p *Payment, calculatedAt, occurredAt time.Time) OverdueAccrued {
	penalty := p.overdue.Penalty
	return OverdueAccrued{
		PaymentID:       p.id.String(),
		UserID:          p.userID.Value().String(),
		DaysOverdue:     p.overdue.DaysOverdue,
		PenaltyAmount:   penalty.Amount().String(),
		PenaltyCurrency: string(penalty.Currency()),
		CalculatedAt:    calculatedAt,
		OccurredAtTime:  occurredAt,
	}
}

func newPaymentPaidEvent(p *Payment, paidAt time.Time) PaymentPaid {
	return PaymentPaid{
		PaymentID:      p.id.String(),
		UserID:         p.userID.Value().String(),
		PaidAt:         paidAt,
		OccurredAtTime: paidAt,
	}
}
