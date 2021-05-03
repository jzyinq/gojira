package gojira

import (
	"testing"
)

func TestFormatTimeSpent(t *testing.T) {
	fixtures := []struct {
		TimeSpentSeconds  int
		expectedTimeSpent string
	}{
		{3600, "1h"},
		{60, "1m"},
		{900, "15m"},
		{7080, "1h58m"},
		{901, "15m"},
		{959, "16m"},
	}

	for _, fixture := range fixtures {
		actualTimeSpent := FormatTimeSpent(fixture.TimeSpentSeconds)
		if actualTimeSpent != fixture.expectedTimeSpent {
			t.Errorf("Incorrect timeSpent - got %s instead of %s", actualTimeSpent, fixture.expectedTimeSpent)
		}
	}
}

func TestCalculateTimeSpent(t *testing.T) {
	fixture := []WorkLog{
		{TimeSpentSeconds: 60},   // 1m
		{TimeSpentSeconds: 3600}, // 1h
		{TimeSpentSeconds: 7200}, // 2h
		{TimeSpentSeconds: 901},  // 15m
		{TimeSpentSeconds: 959},  // 16m
	}
	expectedTimeSpent := "3h32m"

	actualTimeSpent := CalculateTimeSpent(fixture)

	if actualTimeSpent != expectedTimeSpent {
		t.Errorf("Incorrect timeSpent - got %s instead of %s", actualTimeSpent, expectedTimeSpent)
	}
}
