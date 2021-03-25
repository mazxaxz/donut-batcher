package money

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrZeroAmount     = errors.New("provided amount value is zero")
	ErrNegativeAmount = errors.New("provided amount value is negative")
)

func Add(a, b string) (string, error) {
	d1, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}
	d2, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}
	return d1.Add(d2).String(), nil
}

func GreaterThanOrEqual(base, comparer string) (bool, error) {
	d1, err := decimal.NewFromString(base)
	if err != nil {
		return false, err
	}
	d2, err := decimal.NewFromString(comparer)
	if err != nil {
		return false, err
	}
	return d1.GreaterThanOrEqual(d2), nil
}

func CalculateInvestment(amount string) (string, error) {
	value, err := decimal.NewFromString(amount)
	if err != nil {
		return "", err
	}
	if value.IsNegative() {
		return "", ErrNegativeAmount
	}
	if value.IsZero() {
		return "", ErrZeroAmount
	}

	ceiled := value.Ceil()
	investment := ceiled.Sub(value)
	return investment.String(), nil
}
