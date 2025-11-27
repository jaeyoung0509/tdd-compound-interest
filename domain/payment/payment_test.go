package payment

import (
	"testing"
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/domain/user"
	"github.com/stretchr/testify/require"
)

func TestAccrueInterest_CompoundDaily(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	now := base.Add(48 * time.Hour)    // 2 days overdue
	err = p.AccrueInterest(now, 1_000) // 10% daily
	require.NoError(t, err)

	info := p.OverdueInfo()
	require.NotNil(t, info)
	require.Equal(t, StatusOverdue, p.Status())
	require.Equal(t, 2, info.DaysOverdue)

	expectedPenalty := mustKRW(t, 2_100) // day1: +1000, day2: +1100 => 2100
	require.True(t, info.Penalty.Amount().Equal(expectedPenalty.Amount()))
	require.Equal(t, truncateToDate(now), info.CalculatedAt)
}

func TestAccrueInterest_BeforeDueNoOp(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	err = p.AccrueInterest(base, 1_000) // at due date, not overdue yet
	require.NoError(t, err)
	require.Nil(t, p.OverdueInfo())
	require.Equal(t, StatusScheduled, p.Status())
}

func TestAccrueInterest_ContinuesFromExistingPenalty(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	first := base.Add(24 * time.Hour) // 1 day overdue
	require.NoError(t, p.AccrueInterest(first, 1_000))
	require.Equal(t, 1, p.OverdueInfo().DaysOverdue)

	second := base.Add(72 * time.Hour) // 3 days total
	require.NoError(t, p.AccrueInterest(second, 1_000))

	info := p.OverdueInfo()
	require.Equal(t, 3, info.DaysOverdue)
	expectedPenalty := mustKRW(t, 3_310) // day1:1000, day2:1100, day3:1210 => 3310
	require.True(t, info.Penalty.Amount().Equal(expectedPenalty.Amount()))
}

func TestAccrueInterest_FailsWhenPaid(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)
	require.NoError(t, p.Pay(base))

	err = p.AccrueInterest(base.Add(48*time.Hour), 1_000)
	require.ErrorIs(t, err, ErrPaidPaymentCannotOverdue)
	require.Nil(t, p.OverdueInfo())
}

func TestAccrueInterestWith_UsesClockAndRateProvider(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base)
	require.NoError(t, err)

	fakeClock := FixedClock{NowTime: base.Add(48 * time.Hour)} // 2 days
	rate := StaticDailyRate{BPS: 1_000}

	require.NoError(t, p.AccrueInterestWith(fakeClock, rate))

	info := p.OverdueInfo()
	require.NotNil(t, info)
	require.Equal(t, 2, info.DaysOverdue)
	require.Equal(t, truncateToDate(fakeClock.Now()), info.CalculatedAt)
}

func TestPay_DoublePayFails(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base)
	require.NoError(t, err)
	paidAt := base.Add(12 * time.Hour)

	require.NoError(t, p.Pay(paidAt))
	require.ErrorIs(t, p.Pay(paidAt.Add(time.Hour)), ErrPaymentAlreadyPaid)
	require.Equal(t, StatusPaid, p.Status())
}

func TestMarkOverdue_FailsWhenPaid(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)
	p, err := New(uid, amt, base, base)
	require.NoError(t, err)
	require.NoError(t, p.Pay(base.Add(time.Hour)))

	penalty := mustKRW(t, 0)
	err = p.MarkOverdue(base.Add(24*time.Hour), 1, penalty)
	require.ErrorIs(t, err, ErrPaidPaymentCannotOverdue)
}

func TestAccrueInterest_FailsWithZeroNow(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	err = p.AccrueInterest(time.Time{}, 1_000)
	require.ErrorIs(t, err, ErrInvalidOverdueArgs)
	require.Nil(t, p.OverdueInfo())
	require.Equal(t, StatusScheduled, p.Status())
}

func TestAccrueInterest_SameDayDoesNotReaccumulate(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	day1 := base.Add(24 * time.Hour)
	require.NoError(t, p.AccrueInterest(day1, 1_000))
	first := p.OverdueInfo()
	require.NotNil(t, first)

	// Second call within the same calendar day should be a no-op.
	require.NoError(t, p.AccrueInterest(day1.Add(6*time.Hour), 1_000))
	second := p.OverdueInfo()
	require.Equal(t, first.ID, second.ID)
	require.Equal(t, first.DaysOverdue, second.DaysOverdue)
	require.True(t, first.Penalty.Amount().Equal(second.Penalty.Amount()))
	require.Equal(t, first.CalculatedAt, second.CalculatedAt)
}

func TestAccrueInterest_ZeroDailyRateStillMarksOverdue(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)

	now := base.Add(48 * time.Hour)
	require.NoError(t, p.AccrueInterest(now, 0))

	info := p.OverdueInfo()
	require.NotNil(t, info)
	require.Equal(t, StatusOverdue, p.Status())
	require.Equal(t, 2, info.DaysOverdue)
	require.True(t, info.Penalty.IsZero())
}

func TestDaysBetween_HandlesDSTByUsingUTC(t *testing.T) {
	// Simulate DST transition by using a location with DST and ensuring calculation is UTC-based.
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	start := time.Date(2024, 3, 9, 0, 0, 0, 0, loc) // before DST jump
	end := time.Date(2024, 3, 11, 0, 0, 0, 0, loc)  // after DST jump (23h day in between)

	require.Equal(t, 2, daysBetween(start, end))
}

func TestMarkOverdue_ValidationAndDoubleCall(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	uid := mustUserID(t, base)
	amt := mustKRW(t, 10_000)

	p, err := New(uid, amt, base, base.Add(-time.Hour))
	require.NoError(t, err)
	zeroPenalty := mustKRW(t, 0)

	require.ErrorIs(t, p.MarkOverdue(time.Time{}, 1, zeroPenalty), ErrInvalidOverdueArgs)
	require.ErrorIs(t, p.MarkOverdue(base.Add(24*time.Hour), 0, zeroPenalty), ErrInvalidOverdueArgs)

	firstCalc := base.Add(24 * time.Hour)
	require.NoError(t, p.MarkOverdue(firstCalc, 1, zeroPenalty))
	require.Equal(t, StatusOverdue, p.Status())
	first := p.OverdueInfo()
	require.NotNil(t, first)

	require.ErrorIs(t, p.MarkOverdue(firstCalc.Add(24*time.Hour), 1, zeroPenalty), ErrPaymentAlreadyOverdue)
	second := p.OverdueInfo()
	require.Equal(t, first.ID, second.ID)
	require.Equal(t, first.CalculatedAt, second.CalculatedAt)
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
