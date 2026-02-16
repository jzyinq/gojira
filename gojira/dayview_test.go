package gojira

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateTimeSpentInput(t *testing.T) {
	t.Run("accepts empty input", func(t *testing.T) {
		assert.True(t, validateTimeSpentInput("", '1'))
	})

	t.Run("accepts valid digit", func(t *testing.T) {
		assert.True(t, validateTimeSpentInput("1", '1'))
		assert.True(t, validateTimeSpentInput("12", '2'))
	})

	t.Run("accepts 'h' character", func(t *testing.T) {
		assert.True(t, validateTimeSpentInput("1h", 'h'))
	})

	t.Run("accepts 'm' character", func(t *testing.T) {
		assert.True(t, validateTimeSpentInput("30m", 'm'))
	})

	t.Run("accepts space character", func(t *testing.T) {
		assert.True(t, validateTimeSpentInput("1h ", ' '))
	})

	t.Run("rejects invalid characters", func(t *testing.T) {
		assert.False(t, validateTimeSpentInput("1", 'a'))
		assert.False(t, validateTimeSpentInput("1h", 'x'))
		assert.False(t, validateTimeSpentInput("30", '!'))
	})
}

func TestParseDateRange(t *testing.T) {
	t.Run("parses single date", func(t *testing.T) {
		result, err := ParseDateRange("2024-03-15")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), result.StartDate)
		assert.Equal(t, time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), result.EndDate)
		assert.Equal(t, 0, result.NumberOfDays)
	})

	t.Run("parses date range", func(t *testing.T) {
		result, err := ParseDateRange("2024-03-15->2024-03-20")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), result.StartDate)
		assert.Equal(t, time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC), result.EndDate)
		assert.Equal(t, 5, result.NumberOfDays)
	})

	t.Run("returns error for invalid single date", func(t *testing.T) {
		_, err := ParseDateRange("not-a-date")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing date")
	})

	t.Run("returns error for invalid start date in range", func(t *testing.T) {
		_, err := ParseDateRange("bad-date->2024-03-20")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing start date")
	})

	t.Run("returns error for invalid end date in range", func(t *testing.T) {
		_, err := ParseDateRange("2024-03-15->bad-date")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing end date")
	})

	t.Run("handles same start and end date in range", func(t *testing.T) {
		result, err := ParseDateRange("2024-03-15->2024-03-15")
		assert.NoError(t, err)
		assert.Equal(t, result.StartDate, result.EndDate)
		assert.Equal(t, 0, result.NumberOfDays)
	})

	t.Run("handles month-spanning range", func(t *testing.T) {
		result, err := ParseDateRange("2024-01-28->2024-02-03")
		assert.NoError(t, err)
		assert.Equal(t, 6, result.NumberOfDays)
	})
}
