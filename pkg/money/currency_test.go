package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyFrom(t *testing.T) {
	tests := []struct {
		give         string
		wantCurrency Currency
		wantError    error
	}{
		{
			give:         "",
			wantCurrency: "",
			wantError:    ErrInvalidCurrencyCode,
		},
		{
			give:         "US",
			wantCurrency: "",
			wantError:    ErrInvalidCurrencyCode,
		},
		{
			give:         "usd",
			wantCurrency: Currency("USD"),
			wantError:    nil,
		},
		{
			give:         "USD",
			wantCurrency: Currency("USD"),
			wantError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			result, err := CurrencyFrom(tt.give)
			assert.Equal(t, tt.wantCurrency, result)
			assert.Equal(t, tt.wantError, err)
		})
	}
}
