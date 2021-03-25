package money

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		giveA     string
		giveB     string
		want      string
		wantError bool
	}{
		{
			giveA:     "0.01",
			giveB:     "0.001",
			want:      "0.011",
			wantError: false,
		},
		{
			giveA:     "-1.11",
			giveB:     "1.22",
			want:      "0.11",
			wantError: false,
		},
		{
			giveA:     "111.222",
			giveB:     "111.111",
			want:      "222.333",
			wantError: false,
		},
		{
			giveA:     "1.000000000001",
			giveB:     "1.000000000001",
			want:      "2.000000000002",
			wantError: false,
		},
		{
			giveA:     "x",
			giveB:     "x",
			want:      "",
			wantError: true,
		},
		{
			giveA:     "x",
			giveB:     "2.13",
			want:      "",
			wantError: true,
		},
		{
			giveA:     "2.13",
			giveB:     "x",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s + %s", tt.giveA, tt.giveB), func(t *testing.T) {
			result, err := Add(tt.giveA, tt.giveB)
			assert.Equal(t, tt.want, result)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		giveA     string
		giveB     string
		want      bool
		wantError bool
	}{
		{
			giveA:     "111.11",
			giveB:     "100.0",
			want:      true,
			wantError: false,
		},
		{
			giveA:     "99.99",
			giveB:     "100",
			want:      false,
			wantError: false,
		},
		{
			giveA:     "x",
			giveB:     "100",
			want:      false,
			wantError: true,
		},
		{
			giveA:     "99.99",
			giveB:     "x",
			want:      false,
			wantError: true,
		},
		{
			giveA:     "99.99",
			giveB:     "99.99",
			want:      true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s >= %s", tt.giveA, tt.giveB), func(t *testing.T) {
			result, err := GreaterThanOrEqual(tt.giveA, tt.giveB)
			assert.Equal(t, tt.want, result)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculateInvestment(t *testing.T) {
	tests := []struct {
		give      string
		want      string
		wantError error
	}{
		{
			give:      "3.67",
			want:      "0.33",
			wantError: nil,
		},
		{
			give:      "1.999999999",
			want:      "0.000000001",
			wantError: nil,
		},
		{
			give:      "0",
			want:      "",
			wantError: ErrZeroAmount,
		},
		{
			give:      "-2.13",
			want:      "",
			wantError: ErrNegativeAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			result, err := CalculateInvestment(tt.give)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantError, err)
		})
	}
}
