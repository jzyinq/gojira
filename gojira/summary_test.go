package gojira

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkingHoursInMonthToPresentDay(t *testing.T) {
	t.Run("counts weekdays excluding holidays", func(t *testing.T) {
		holidays := &Holidays{
			{Date: "2024-01-01", Name: "New Year"},
		}
		// January 2024: 23 weekdays, minus 1 holiday = 22 working days = 176 hours
		hours := workingHoursInMonthToPresentDay(2024, time.January, holidays)
		assert.Equal(t, 176, hours)
	})

	t.Run("counts weekdays with no holidays", func(t *testing.T) {
		holidays := &Holidays{}
		// January 2024: 23 weekdays, no holidays = 23 * 8 = 184 hours
		hours := workingHoursInMonthToPresentDay(2024, time.January, holidays)
		assert.Equal(t, 184, hours)
	})

	t.Run("counts weekdays with multiple holidays", func(t *testing.T) {
		holidays := &Holidays{
			{Date: "2024-01-01", Name: "New Year"},
			{Date: "2024-01-15", Name: "MLK Day"},
		}
		// January 2024: 23 weekdays, minus 2 holidays = 21 * 8 = 168 hours
		hours := workingHoursInMonthToPresentDay(2024, time.January, holidays)
		assert.Equal(t, 168, hours)
	})

	t.Run("weekend holiday does not reduce count", func(t *testing.T) {
		holidays := &Holidays{
			// January 6, 2024 is a Saturday
			{Date: "2024-01-06", Name: "Weekend Holiday"},
		}
		// January 2024: 23 weekdays, holiday falls on Saturday so no reduction
		hours := workingHoursInMonthToPresentDay(2024, time.January, holidays)
		assert.Equal(t, 184, hours)
	})

	t.Run("february leap year", func(t *testing.T) {
		holidays := &Holidays{}
		// February 2024 (leap year): 29 days, 21 weekdays = 168 hours
		hours := workingHoursInMonthToPresentDay(2024, time.February, holidays)
		assert.Equal(t, 168, hours)
	})
}

func TestWorkingHoursAbsoluteDiff(t *testing.T) {
	t.Run("positive difference (under-logged)", func(t *testing.T) {
		// 160 working hours expected, 150 hours of time spent in seconds
		diff := workingHoursAbsoluteDiff(160, 150*3600)
		assert.Equal(t, 10*3600, diff) // 10 hours difference
	})

	t.Run("negative difference (over-logged)", func(t *testing.T) {
		// 160 working hours expected, 170 hours of time spent in seconds
		diff := workingHoursAbsoluteDiff(160, 170*3600)
		assert.Equal(t, 10*3600, diff) // absolute value: still 10 hours
	})

	t.Run("exact match returns zero", func(t *testing.T) {
		diff := workingHoursAbsoluteDiff(160, 160*3600)
		assert.Equal(t, 0, diff)
	})

	t.Run("zero working hours with some time spent", func(t *testing.T) {
		diff := workingHoursAbsoluteDiff(0, 3600)
		assert.Equal(t, 3600, diff) // 1 hour
	})

	t.Run("zero everything", func(t *testing.T) {
		diff := workingHoursAbsoluteDiff(0, 0)
		assert.Equal(t, 0, diff)
	})
}
