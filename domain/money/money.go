package money

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidCurrency  = errors.New("invalid currency")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

// Currency expects an ISO-like currency identifier.
type Currency string

const (
	CurrencyKRW Currency = "KRW"
	CurrencyUSD Currency = "USD"
)

func currencyScale(c Currency) (int32, error) {
	switch c {
	case CurrencyKRW:
		return 0, nil
	case CurrencyUSD:
		return 2, nil
	default:
		return 0, ErrInvalidCurrency
	}
}

// Money represents an amount normalized to the configured currency scale.
type Money struct {
	amount   decimal.Decimal
	currency Currency
	scale    int32
}

// New는 금액을 통화 스케일에 맞춰 반올림하여 Money를 생성합니다.
func New(amount decimal.Decimal, currency Currency) (Money, error) {
	scale, err := currencyScale(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{
		amount:   amount.Round(scale),
		currency: currency,
		scale:    scale,
	}, nil
}

// FromMinor builds Money from minor units (e.g., KRW won, USD cents).
func FromMinor(minor int64, currency Currency) (Money, error) {
	scale, err := currencyScale(currency)
	if err != nil {
		return Money{}, err
	}
	dec := decimal.NewFromInt(minor).Shift(-scale)
	return Money{
		amount:   dec,
		currency: currency,
		scale:    scale,
	}, nil
}

// Zero returns a zero-valued Money for the given currency.
func Zero(currency Currency) (Money, error) {
	scale, err := currencyScale(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{
		amount:   decimal.Zero,
		currency: currency,
		scale:    scale,
	}, nil
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) Amount() decimal.Decimal {
	return m.amount
}

func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	res := m.amount.Add(other.amount).Round(m.scale)
	return Money{amount: res, currency: m.currency, scale: m.scale}, nil
}

func (m Money) Sub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	res := m.amount.Sub(other.amount).Round(m.scale)
	return Money{amount: res, currency: m.currency, scale: m.scale}, nil
}

// MulBPS는 basis points(만분율)로 배수를 적용합니다.
// 100bps = 1%, 10_000bps = 100%. 반올림은 통화 스케일에서 half-up.
func (m Money) MulBPS(bps int64) Money {
	if bps == 0 || m.amount.IsZero() {
		return Money{currency: m.currency, scale: m.scale}
	}
	rate := decimal.NewFromInt(bps).Div(decimal.NewFromInt(10_000))
	delta := m.amount.Mul(rate).Round(m.scale)
	return Money{amount: delta, currency: m.currency, scale: m.scale}
}

// ApplyBPS는 원금에 bps 이율을 적용한 금액을 반환합니다.
func (m Money) ApplyBPS(bps int64) (Money, error) {
	delta := m.MulBPS(bps)
	return m.Add(delta)
}

func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.amount.StringFixed(m.scale), m.currency)
}
