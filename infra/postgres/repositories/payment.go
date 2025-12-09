package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/domain/user"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/sqlc/generated"
	"github.com/shopspring/decimal"
)

type PaymentRepository struct {
	queries *generated.Queries
}

func NewPaymentRepository(db generated.DBTX) *PaymentRepository {
	return &PaymentRepository{queries: generated.New(db)}
}

func (r *PaymentRepository) Get(ctx context.Context, id shared.ID) (*payment.Payment, error) {
	row, err := r.queries.GetPayment(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, payment.ErrPaymentNotFound
		}
		return nil, err
	}

	// Get latest overdue info
	var overdueInfo *payment.OverdueInfo
	overdueRow, err := r.queries.GetLatestOverdue(ctx, id.String())
	if err == nil {
		overdueInfo = toOverdueDomain(overdueRow)
	}

	return toPaymentDomain(row, overdueInfo)
}

func (r *PaymentRepository) Save(ctx context.Context, p *payment.Payment) error {
	params := fromPaymentDomain(p)

	if err := r.queries.UpsertPayment(ctx, params); err != nil {
		return err
	}

	// Save overdue info if present
	if p.OverdueInfo() != nil {
		overdueParams := fromOverdueDomain(p.ID(), p.OverdueInfo())
		if err := r.queries.InsertOverdue(ctx, overdueParams); err != nil {
			return err
		}
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Mapper functions: toDomain / fromDomain
// ──────────────────────────────────────────────────────────────────────────────

func toPaymentDomain(row generated.Payment, overdue *payment.OverdueInfo) (*payment.Payment, error) {
	id, err := shared.ParseID(row.ID)
	if err != nil {
		return nil, err
	}

	userID, err := shared.ParseID(row.UserID)
	if err != nil {
		return nil, err
	}

	amount, err := toMoneyDomain(row.Amount, row.Currency)
	if err != nil {
		return nil, err
	}

	var paidAt *time.Time
	if row.PaidAt.Valid {
		paidAt = &row.PaidAt.Time
	}

	return payment.Reconstitute(
		id,
		user.IDFromShared(userID),
		amount,
		row.DueDate.Time,
		paidAt,
		payment.Status(row.Status),
		overdue,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}

func fromPaymentDomain(p *payment.Payment) generated.UpsertPaymentParams {
	var paidAt pgtype.Timestamptz
	if p.PaidAt() != nil {
		paidAt = pgtype.Timestamptz{Time: *p.PaidAt(), Valid: true}
	}

	return generated.UpsertPaymentParams{
		ID:        p.ID().String(),
		UserID:    p.UserID().Value().String(),
		Amount:    toNumeric(p.Amount().Amount()),
		Currency:  string(p.Amount().Currency()),
		DueDate:   pgtype.Date{Time: p.DueDate(), Valid: true},
		PaidAt:    paidAt,
		Status:    string(p.Status()),
		CreatedAt: pgtype.Timestamptz{Time: p.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: p.UpdatedAt(), Valid: true},
	}
}

func toOverdueDomain(row generated.PaymentOverdue) *payment.OverdueInfo {
	id, _ := shared.ParseID(row.ID)
	penalty, _ := toMoneyDomain(row.Penalty, row.PenaltyCurrency)

	return &payment.OverdueInfo{
		ID:           id,
		IsOverdue:    row.IsOverdue,
		DaysOverdue:  int(row.DaysOverdue),
		Penalty:      penalty,
		CalculatedAt: row.CalculatedAt.Time,
	}
}

func fromOverdueDomain(paymentID shared.ID, o *payment.OverdueInfo) generated.InsertOverdueParams {
	return generated.InsertOverdueParams{
		ID:              o.ID.String(),
		PaymentID:       paymentID.String(),
		IsOverdue:       o.IsOverdue,
		DaysOverdue:     int32(o.DaysOverdue),
		Penalty:         toNumeric(o.Penalty.Amount()),
		PenaltyCurrency: string(o.Penalty.Currency()),
		CalculatedAt:    pgtype.Date{Time: o.CalculatedAt, Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Money conversion helpers
// ──────────────────────────────────────────────────────────────────────────────

func toMoneyDomain(n pgtype.Numeric, currency string) (money.Money, error) {
	if !n.Valid {
		return money.Zero(money.Currency(currency))
	}

	// Convert pgtype.Numeric to decimal.Decimal
	var d decimal.Decimal
	if n.Int != nil {
		d = decimal.NewFromBigInt(n.Int, n.Exp)
	}

	return money.New(d, money.Currency(currency))
}

func toNumeric(d decimal.Decimal) pgtype.Numeric {
	// Convert decimal.Decimal to pgtype.Numeric
	var n pgtype.Numeric
	_ = n.Scan(d.String())
	return n
}

// Ensure PaymentRepository implements payment.Repository interface
var _ payment.Repository = (*PaymentRepository)(nil)
