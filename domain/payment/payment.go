package payment

import (
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/domain/user"
)

// Payment is the aggregate root that encapsulates payment lifecycle transitions.
type Payment struct {
	id        shared.ID
	userID    user.ID
	amount    money.Money
	dueDate   time.Time
	paidAt    *time.Time
	status    Status
	overdue   *OverdueInfo
	createdAt time.Time
	updatedAt time.Time
}

func New(userID user.ID, amount money.Money, dueDate time.Time, now time.Time) (*Payment, error) {
	if userID.IsZero() {
		return nil, ErrInvalidUserID
	}
	if dueDate.IsZero() {
		return nil, ErrInvalidDueDate
	}
	if now.IsZero() {
		now = time.Now()
	}

	return &Payment{
		id:        shared.NewID(),
		userID:    userID,
		amount:    amount,
		dueDate:   dueDate,
		status:    StatusScheduled,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (p *Payment) ID() shared.ID {
	return p.id
}

func (p *Payment) UserID() user.ID {
	return p.userID
}

func (p *Payment) Amount() money.Money {
	return p.amount
}

func (p *Payment) DueDate() time.Time {
	return p.dueDate
}

func (p *Payment) PaidAt() *time.Time {
	if p.paidAt == nil {
		return nil
	}
	t := *p.paidAt
	return &t
}

func (p *Payment) Status() Status {
	return p.status
}

func (p *Payment) OverdueInfo() *OverdueInfo {
	if p.overdue == nil {
		return nil
	}
	info := *p.overdue
	return &info
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Payment) Pay(paidAt time.Time) error {
	if paidAt.IsZero() {
		return ErrInvalidPaidAt
	}
	if p.status == StatusPaid {
		return ErrPaymentAlreadyPaid
	}

	p.paidAt = &paidAt
	p.status = StatusPaid
	p.updatedAt = paidAt
	return nil
}

func (p *Payment) MarkOverdue(calculatedAt time.Time, daysOverdue int, penalty money.Money) error {
	if calculatedAt.IsZero() || daysOverdue <= 0 {
		return ErrInvalidOverdueArgs
	}
	if p.status == StatusPaid {
		return ErrPaidPaymentCannotOverdue
	}
	if p.status == StatusOverdue {
		return ErrPaymentAlreadyOverdue
	}

	p.overdue = &OverdueInfo{
		ID:           shared.NewID(),
		IsOverdue:    true,
		DaysOverdue:  daysOverdue,
		Penalty:      penalty,
		CalculatedAt: calculatedAt,
	}
	p.status = StatusOverdue
	p.updatedAt = calculatedAt
	return nil
}
