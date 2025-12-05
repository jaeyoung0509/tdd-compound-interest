package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/money"
	dp "github.com/jaeyoung0509/compound-interest/domain/payment"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/domain/user"
	"github.com/stretchr/testify/require"
)

func TestAccruePayment_Success(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := dp.New(uid, amt, base, base)
	require.NoError(t, err)

	repo := NewInMemoryPaymentRepo()
	repo.Seed(p)

	svc := NewService(repo, dp.FixedClock{NowTime: base.Add(48 * time.Hour)}, dp.StaticDailyRate{BPS: 1_000})

	updated, err := svc.AccruePayment(context.Background(), p.ID())
	require.NoError(t, err)
	require.Equal(t, 1, repo.SaveCount())

	info := updated.OverdueInfo()
	require.NotNil(t, info)
	require.Equal(t, dp.StatusOverdue, updated.Status())
	require.Equal(t, 2, info.DaysOverdue)

	expectedPenalty := mustKRW(t, 2_100)
	require.True(t, info.Penalty.Amount().Equal(expectedPenalty.Amount()))
}

func TestAccruePayment_PaidSkipsSave(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := dp.New(uid, amt, base, base)
	require.NoError(t, err)
	require.NoError(t, p.Pay(base.Add(time.Hour)))

	repo := NewInMemoryPaymentRepo()
	repo.Seed(p)

	svc := NewService(repo, dp.FixedClock{NowTime: base.Add(48 * time.Hour)}, dp.StaticDailyRate{BPS: 1_000})

	_, err = svc.AccruePayment(context.Background(), p.ID())
	require.ErrorIs(t, err, dp.ErrPaidPaymentCannotOverdue)
	require.Equal(t, 0, repo.SaveCount())
}

func TestAccruePayment_NotFound(t *testing.T) {
	repo := NewInMemoryPaymentRepo()
	svc := NewService(repo, dp.FixedClock{NowTime: time.Now()}, dp.StaticDailyRate{BPS: 1_000})

	_, err := svc.AccruePayment(context.Background(), shared.NewID())
	require.ErrorIs(t, err, dp.ErrPaymentNotFound)
	require.Equal(t, 0, repo.SaveCount())
}

func TestAccruePayment_SaveError(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := dp.New(uid, amt, base, base)
	require.NoError(t, err)

	repo := NewInMemoryPaymentRepo()
	repo.Seed(p)
	repo.SaveErr = errors.New("save fail")

	svc := NewService(repo, dp.FixedClock{NowTime: base.Add(48 * time.Hour)}, dp.StaticDailyRate{BPS: 1_000})

	_, err = svc.AccruePayment(context.Background(), p.ID())
	require.Error(t, err)
	require.EqualError(t, err, "save fail")
	require.Equal(t, 0, repo.SaveCount())
}

func mustUserID(t *testing.T, now time.Time) user.ID {
	u, err := user.New("tester", now)
	require.NoError(t, err)
	return u.ID()
}

func mustKRW(t *testing.T, minor int64) money.Money {
	m, err := money.FromMinor(minor, money.CurrencyKRW)
	require.NoError(t, err)
	return m
}
