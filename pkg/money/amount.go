package money

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrZeroAmount     = errors.New("provided amount value is zero")
	ErrNegativeAmount = errors.New("provided amount value is negative")
)

func CalculateInvestment(amount string) (decimal.Decimal, error) {
	value, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Decimal{}, err
	}
	if value.IsNegative() {
		return decimal.Decimal{}, ErrNegativeAmount
	}
	if value.IsZero() {
		return decimal.Decimal{}, ErrZeroAmount
	}

	ceiled := value.Ceil()
	investment := ceiled.Sub(value)
	return investment, nil
}
