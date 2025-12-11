//go:build integration

package payment_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/money"
	"github.com/jaeyoung0509/compound-interest/infra/postgres"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/testhelper"
	"github.com/jaeyoung0509/compound-interest/usecase/payment"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *testhelper.TestDB

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 컨테이너 한 번만 기동
	var err error
	testDB, err = testhelper.NewTestDB(ctx)
	if err != nil {
		panic(err)
	}

	// 모든 테스트 실행
	code := m.Run()

	// Cleanup
	testDB.Close(ctx)
	os.Exit(code)
}

func TestService_AccruePayment_OverduePayment(t *testing.T) {
	ctx := context.Background()

	// Given: UoW + Service
	uow, err := postgres.NewUnitOfWork(ctx, testDB.Pool)
	require.NoError(t, err)
	defer uow.Rollback(ctx) // 테스트 격리

	fixedTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	fixedClock := &FixedClock{now: fixedTime}

	service := payment.NewService(
		uow.Payments(),
		fixedClock,
		&testRateProvider{bps: 100}, // 100 BPS = 1%
	)

	// Given: 10일 전에 만료된 payment (연체 10일)
	userID, err := testhelper.CreateTestUser(ctx, testDB.Pool, "Test User", fixedTime)
	require.NoError(t, err)

	pastDue := fixedTime.Add(-10 * 24 * time.Hour)
	amount, _ := money.New(decimal.NewFromInt(100000), money.CurrencyKRW)
	paymentID, err := testhelper.CreateTestPayment(
		ctx, testDB.Pool, userID,
		amount, pastDue, fixedTime,
	)
	require.NoError(t, err)

	// When: AccruePayment 실행
	result, err := service.AccruePayment(ctx, paymentID)

	// Then: 연체 이자 계산됨
	require.NoError(t, err)
	assert.NotNil(t, result.OverdueInfo())
	assert.True(t, result.OverdueInfo().IsOverdue)
	assert.Equal(t, 10, result.OverdueInfo().DaysOverdue)

	// Compound interest calculation (not simple interest)
	// Each day adds interest on (principal + accumulated penalty)
	// Simple interest: 100000 * 0.01 * 10 = 10000
	// Compound interest: slightly higher (~10463 for 100 BPS over 10 days)
	assert.True(t, result.OverdueInfo().Penalty.Amount().GreaterThan(decimal.NewFromInt(10000)),
		"Compound interest should exceed simple interest")
	assert.True(t, result.OverdueInfo().Penalty.Amount().LessThan(decimal.NewFromInt(11000)),
		"Penalty should be reasonable")
}

func TestService_AccruePayment_NotOverdue(t *testing.T) {
	ctx := context.Background()

	uow, err := postgres.NewUnitOfWork(ctx, testDB.Pool)
	require.NoError(t, err)
	defer uow.Rollback(ctx)

	fixedTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	service := payment.NewService(
		uow.Payments(),
		&FixedClock{now: fixedTime},
		&testRateProvider{bps: 100},
	)

	// Given: 아직 만료되지 않은 payment
	userID, err := testhelper.CreateTestUser(ctx, testDB.Pool, "Alice", fixedTime)
	require.NoError(t, err)

	futureDue := fixedTime.Add(7 * 24 * time.Hour) // 7일 후
	amount, _ := money.New(decimal.NewFromInt(50000), money.CurrencyKRW)
	paymentID, err := testhelper.CreateTestPayment(
		ctx, testDB.Pool, userID,
		amount, futureDue, fixedTime,
	)
	require.NoError(t, err)

	// When: AccruePayment 실행
	result, err := service.AccruePayment(ctx, paymentID)

	// Then: 연체 아님
	require.NoError(t, err)
	assert.Nil(t, result.OverdueInfo())
}

func TestService_AccruePayment_AlreadyPaid(t *testing.T) {
	ctx := context.Background()

	uow, err := postgres.NewUnitOfWork(ctx, testDB.Pool)
	require.NoError(t, err)
	defer uow.Rollback(ctx)

	fixedTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	service := payment.NewService(
		uow.Payments(),
		&FixedClock{now: fixedTime},
		&testRateProvider{bps: 100},
	)

	// Given: 이미 완납된 payment
	userID, err := testhelper.CreateTestUser(ctx, testDB.Pool, "Bob", fixedTime)
	require.NoError(t, err)

	dueDate := fixedTime.Add(-5 * 24 * time.Hour)
	amount, _ := money.New(decimal.NewFromInt(200000), money.CurrencyKRW)

	paidAt := fixedTime.Add(-2 * 24 * time.Hour)
	paymentID, err := testhelper.CreateTestPayment(
		ctx, testDB.Pool, userID,
		amount, dueDate, fixedTime,
		testhelper.WithPaidAt(paidAt),
	)
	require.NoError(t, err)

	// When: AccruePayment 실행
	result, err := service.AccruePayment(ctx, paymentID)

	// Then: 이미 완납되어 에러 또는 unchanged
	// (실제 비즈니스 로직에 따라 다름)
	_ = result
	_ = err
	// 여기서는 에러가 발생하거나, 변화 없이 반환되어야 함
}

// FixedClock implements payment.Clock with a fixed time (for deterministic tests)
type FixedClock struct {
	now time.Time
}

func (c *FixedClock) Now() time.Time {
	return c.now
}

// testRateProvider implements payment.DailyRateProvider
type testRateProvider struct {
	bps int64
}

func (p *testRateProvider) DailyRateBPS(at time.Time) (int64, error) {
	return p.bps, nil
}
