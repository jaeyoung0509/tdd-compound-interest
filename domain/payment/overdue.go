package payment

import (
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

type OverdueInfo struct {
	ID           shared.ID
	IsOverdue    bool
	DaysOverdue  int
	Penalty      money.Money
	CalculatedAt time.Time
}
