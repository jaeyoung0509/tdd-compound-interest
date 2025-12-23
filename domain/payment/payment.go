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

const maxOverdueDays = 365*3 + 1 // three years with a leap-day allowance

func New(userID user.ID, amount money.Money, dueDate time.Time, now time.Time) (*Payment, error) {
	if userID.IsZero() {
		return nil, ErrInvalidUserID
	}
	if amount.Amount().Sign() <= 0 {
		return nil, ErrInvalidAmount
	}
	if dueDate.IsZero() {
		return nil, ErrInvalidDueDate
	}
	if now.IsZero() {
		now = time.Now()
	}
	dueDate = truncateToDate(dueDate)
	if truncateToDate(now).After(dueDate) {
		return nil, ErrDueDateInPast
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
	if truncateToDate(paidAt).Before(truncateToDate(p.dueDate)) {
		return ErrPaidBeforeDueDate
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

// AccrueInterestWith pulls time and rate from collaborators to simplify wiring.
func (p *Payment) AccrueInterestWith(clock Clock, rateProvider DailyRateProvider) error {
	now := clock.Now()
	rate, err := rateProvider.DailyRateBPS(now)
	if err != nil {
		return err
	}
	return p.AccrueInterest(now, rate)
}

// AccrueInterest compounds daily interest from the due date (or last accrual)
// up to the provided time using basis points per day. No-op if not past due.
func (p *Payment) AccrueInterest(now time.Time, dailyRateBPS int64) error {
	if now.IsZero() {
		return ErrInvalidOverdueArgs
	}
	if p.status == StatusPaid {
		return ErrPaidPaymentCannotOverdue
	}

	anchor := p.dueDate
	penalty, err := money.Zero(p.amount.Currency())
	if err != nil {
		return err
	}
	accumulatedDays := 0
	if p.overdue != nil {
		anchor = p.overdue.CalculatedAt
		penalty = p.overdue.Penalty
		accumulatedDays = p.overdue.DaysOverdue
	}

	if !now.After(anchor) {
		return nil
	}

	days := daysBetween(anchor, now)
	if days <= 0 {
		return nil
	}

	totalDays := accumulatedDays + days
	if totalDays > maxOverdueDays {
		return ErrOverduePeriodTooLong
	}

	currentPenalty := penalty
	for i := 0; i < days; i++ {
		base, err := p.amount.Add(currentPenalty)
		if err != nil {
			return err
		}
		delta := base.MulBPS(dailyRateBPS)
		currentPenalty, err = currentPenalty.Add(delta)
		if err != nil {
			return err
		}
	}

	p.overdue = &OverdueInfo{
		ID:           shared.NewID(),
		IsOverdue:    true,
		DaysOverdue:  totalDays,
		Penalty:      currentPenalty,
		CalculatedAt: truncateToDate(now),
	}
	p.status = StatusOverdue
	p.updatedAt = now
	return nil
}

func truncateToDate(t time.Time) time.Time {
	y, m, d := t.In(time.UTC).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func daysBetween(start, end time.Time) int {
	s := truncateToDate(start)
	e := truncateToDate(end)
	if !e.After(s) {
		return 0
	}
	return int(e.Sub(s).Hours() / 24)
}

// Reconstitute rebuilds a Payment from persistence layer without validation.
// This should only be used by the repository layer.
func Reconstitute(
	id shared.ID,
	userID user.ID,
	amount money.Money,
	dueDate time.Time,
	paidAt *time.Time,
	status Status,
	overdue *OverdueInfo,
	createdAt, updatedAt time.Time,
) *Payment {
	return &Payment{
		id:        id,
		userID:    userID,
		amount:    amount,
		dueDate:   dueDate,
		paidAt:    paidAt,
		status:    status,
		overdue:   overdue,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
