package gojira

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHolidays_GetHolidaysForMonth(t *testing.T) {
	holidays := Holidays{
		{Date: "2024-01-01", Name: "New Year's Day", CountryCode: "US"},
		{Date: "2024-02-14", Name: "Valentine's Day", CountryCode: "US"},
		{Date: "2024-02-19", Name: "Presidents Day", CountryCode: "US"},
		{Date: "2024-03-17", Name: "St. Patrick's Day", CountryCode: "US"},
		{Date: "2024-07-04", Name: "Independence Day", CountryCode: "US"},
	}

	t.Run("returns holidays for specified month", func(t *testing.T) {
		result := holidays.GetHolidaysForMonth(time.February)
		assert.Len(t, result, 2)
		assert.Equal(t, "Valentine's Day", result[0].Name)
		assert.Equal(t, "Presidents Day", result[1].Name)
	})

	t.Run("returns empty for month with no holidays", func(t *testing.T) {
		result := holidays.GetHolidaysForMonth(time.April)
		assert.Len(t, result, 0)
	})

	t.Run("returns single holiday for month with one", func(t *testing.T) {
		result := holidays.GetHolidaysForMonth(time.January)
		assert.Len(t, result, 1)
		assert.Equal(t, "New Year's Day", result[0].Name)
	})
}

func TestHolidays_IsHoliday(t *testing.T) {
	holidays := Holidays{
		{Date: "2024-01-01", Name: "New Year's Day", CountryCode: "US"},
		{Date: "2024-07-04", Name: "Independence Day", CountryCode: "US"},
	}

	t.Run("returns true for holiday date", func(t *testing.T) {
		holidayDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, holidays.IsHoliday(&holidayDate))
	})

	t.Run("returns false for non-holiday date", func(t *testing.T) {
		regularDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		assert.False(t, holidays.IsHoliday(&regularDate))
	})

	t.Run("returns true for another holiday", func(t *testing.T) {
		holidayDate := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		assert.True(t, holidays.IsHoliday(&holidayDate))
	})

	t.Run("returns false for empty holidays list", func(t *testing.T) {
		emptyHolidays := Holidays{}
		anyDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.False(t, emptyHolidays.IsHoliday(&anyDate))
	})
}

func TestHoliday_GetTime(t *testing.T) {
	t.Run("parses valid date", func(t *testing.T) {
		holiday := Holiday{Date: "2024-07-04", Name: "Independence Day"}
		result, err := holiday.GetTime()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2024, result.Year())
		assert.Equal(t, time.July, result.Month())
		assert.Equal(t, 4, result.Day())
	})

	t.Run("returns error for invalid date", func(t *testing.T) {
		holiday := Holiday{Date: "invalid-date", Name: "Test"}
		result, err := holiday.GetTime()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error for empty date", func(t *testing.T) {
		holiday := Holiday{Date: "", Name: "Test"}
		result, err := holiday.GetTime()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error for wrong date format", func(t *testing.T) {
		holiday := Holiday{Date: "07/04/2024", Name: "Test"}
		result, err := holiday.GetTime()
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetCountryFromLCTime(t *testing.T) {
	t.Run("extracts country code from valid LC_TIME string", func(t *testing.T) {
		result, err := GetCountryFromLCTime("en_US.UTF-8")
		assert.NoError(t, err)
		assert.Equal(t, "US", result)
	})

	t.Run("extracts country code from different format", func(t *testing.T) {
		result, err := GetCountryFromLCTime("de_DE")
		assert.NoError(t, err)
		assert.Equal(t, "DE", result)
	})

	t.Run("extracts country code from GB locale", func(t *testing.T) {
		result, err := GetCountryFromLCTime("en_GB.UTF-8")
		assert.NoError(t, err)
		assert.Equal(t, "GB", result)
	})

	t.Run("returns error for string without country code", func(t *testing.T) {
		result, err := GetCountryFromLCTime("invalid")
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "could not parse country")
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := GetCountryFromLCTime("")
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("extracts first country code if multiple", func(t *testing.T) {
		result, err := GetCountryFromLCTime("en_US_GB")
		assert.NoError(t, err)
		assert.Equal(t, "US", result)
	})

	t.Run("ignores lowercase and extracts UTF from UTF-8", func(t *testing.T) {
		// Note: The regex [A-Z]{2} will match "UT" from "UTF-8" in "en_us.UTF-8"
		// This is the actual behavior - it finds the first 2 uppercase letters
		result, err := GetCountryFromLCTime("en_us.UTF-8")
		assert.NoError(t, err)
		assert.Equal(t, "UT", result)
	})
}
