package money

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNewMoneyValidation(t *testing.T) {
	_, err := New(decimal.NewFromInt(100), "")
	require.ErrorIs(t, err, ErrInvalidCurrency)
}

func TestAddAndSubSameCurrency(t *testing.T) {
	m1, _ := New(decimal.NewFromInt(1000), CurrencyKRW)
	m2, _ := New(decimal.NewFromInt(500), CurrencyKRW)

	sum, err := m1.Add(m2)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(1500), sum.Amount())

	diff, err := m1.Sub(m2)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(500), diff.Amount())
}

func TestCurrencyMismatch(t *testing.T) {
	m1, _ := New(decimal.NewFromInt(1000), CurrencyKRW)
	m2, _ := New(decimal.NewFromInt(500), CurrencyUSD)

	_, err := m1.Add(m2)
	require.ErrorIs(t, err, ErrCurrencyMismatch)

	_, err = m1.Sub(m2)
	require.ErrorIs(t, err, ErrCurrencyMismatch)
}

func TestMulBPSRoundingHalfUp(t *testing.T) {
	m, _ := New(decimal.NewFromInt(1000), CurrencyKRW) // 1000원

	got := m.MulBPS(125)                                   // 1.25%
	require.Equal(t, decimal.NewFromInt(13), got.Amount()) // 12.5 -> 13 half-up

	gotNeg := m.MulBPS(-125) // -1.25%
	require.Equal(t, decimal.NewFromInt(-13), gotNeg.Amount())
}

func TestApplyBPS(t *testing.T) {
	m, _ := New(decimal.NewFromInt(10_000), CurrencyKRW) // 10,000원

	withRate, err := m.ApplyBPS(1000) // +10%
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(11_000), withRate.Amount())
}
