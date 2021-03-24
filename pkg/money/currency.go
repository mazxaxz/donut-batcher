package money

import (
	"errors"
	"strings"
)

// Currency representation in ISO 4217 standard
type Currency string

var (
	ErrInvalidCurrencyCode = errors.New("currency code does not match ISO 4217 standard")
)

func CurrencyFrom(input string) (Currency, error) {
	if l := len(input); l != 3 {
		return "", ErrInvalidCurrencyCode
	}
	return Currency(strings.ToUpper(input)), nil
}
