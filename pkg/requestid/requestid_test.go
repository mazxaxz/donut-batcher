package requestid

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should assign given value to context values", func(t *testing.T) {
		// arrange
		const rid = "request_id"
		parent := context.Background()

		// act
		ctx := New(parent, rid)

		// assert
		result, _ := ctx.Value(contextKeyRequestID).(string)
		assert.Equal(t, rid, result)
	})

	t.Run("should assign generated value", func(t *testing.T) {
		// arrange
		parent := context.Background()

		// act
		ctx := New(parent, "")

		// assert
		result, _ := ctx.Value(contextKeyRequestID).(string)
		assert.NotEqual(t, unableToCorrelate, result)
		assert.NotEqual(t, "", result)
	})
}

func TestNewRequestID(t *testing.T) {
	t.Run("should generate new request id", func(t *testing.T) {
		// arrange

		// act
		result := NewRequestID()

		// assert
		assert.NotEqual(t, unableToCorrelate, result)
		assert.NotEqual(t, "", result)
		assert.True(t, strings.HasPrefix(result, "|:"))
	})
}

func TestFrom(t *testing.T) {
	t.Run("should return existing request id", func(t *testing.T) {
		// arrange
		const rid = "request_id"
		ctx := New(context.Background(), rid)

		// act
		result, exists := From(ctx)

		// assert
		assert.Equal(t, rid, result)
		assert.True(t, exists)
	})

	t.Run("should return unable to correlate request id", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		result, exists := From(ctx)

		// assert
		assert.Equal(t, unableToCorrelate, result)
		assert.False(t, exists)
	})
}
