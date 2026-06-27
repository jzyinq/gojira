package gojira

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestFormatTimeSpent(t *testing.T) {
	fixtures := []struct {
		TimeSpentSeconds  int
		expectedTimeSpent string
	}{
		{3600, "1h"},
		{60, "1m"},
		{900, "15m"},
		{7080, "1h 58m"},
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
	fixture := []*Worklog{
		{TimeSpentSeconds: 60},   // 1m
		{TimeSpentSeconds: 3600}, // 1h
		{TimeSpentSeconds: 7200}, // 2h
		{TimeSpentSeconds: 901},  // 15m
		{TimeSpentSeconds: 959},  // 16m
	}
	expectedTimeSpent := "3h 32m"

	actualTimeSpent := FormatTimeSpent(CalculateTimeSpent(fixture))

	if actualTimeSpent != expectedTimeSpent {
		t.Errorf("Incorrect timeSpent - got %s instead of %s", actualTimeSpent, expectedTimeSpent)
	}
}

func TestFindIssueKeyInString(t *testing.T) {
	fixtures := []string{
		"TICKET-999",
		"https://instance.atlassian.net/secure/RapidBoard.jspa?rapidView=84&projectKey=TICKET&" +
			"view=planning&selectedIssue=TICKET-999&issueLimit=100",
		"https://instance.atlassian.net/browse/TICKET-999",
		"anythingreallyTICKET-999COULDBEHERE",
		"COULDBEHERE_TICKET-999cq334q5c3v",
	}
	expectedIssueKey := "TICKET-999"

	for _, fixture := range fixtures {
		actualIssueKey := FindIssueKeyInString(fixture)
		if actualIssueKey != expectedIssueKey {
			t.Errorf("Incorrect timeSpent - got %s instead of %s", actualIssueKey, expectedIssueKey)
		}
	}
}

func TestTimeSpentToSeconds(t *testing.T) {
	fixtures := []struct {
		TimeSpent                  string
		expectedTimeSpentInSeconds int
	}{
		{"1h 30m", 5400},
		{testDuration1h30m, 5400},
		{"29m", 1740},
	}

	for _, fixture := range fixtures {
		timeSpentInSeconds := TimeSpentToSeconds(fixture.TimeSpent)
		if timeSpentInSeconds != fixture.expectedTimeSpentInSeconds {
			t.Errorf("Incorrect timeSpent - got %d instead of %d", timeSpentInSeconds, fixture.expectedTimeSpentInSeconds)
		}
	}
}

func TestGetTimeSpentColorTag(t *testing.T) {
	hours := 8

	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"no time logged", 0, "[white]"},
		{"under target", 7200, "[orange]"},   // 2h when target is 8h
		{"exactly target", 28800, "[green]"}, // 8h
		{"over target", 32400, "[blue]"},     // 9h
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimeSpentColorTag(tt.seconds, hours)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTimeSpentColor(t *testing.T) {
	hours := 8

	tests := []struct {
		name     string
		seconds  int
		expected tcell.Color
	}{
		{"no time logged", 0, tcell.ColorWhite},
		{"under target", 7200, tcell.ColorOrange},   // 2h when target is 8h
		{"exactly target", 28800, tcell.ColorGreen}, // 8h
		{"over target", 32400, tcell.ColorBlue},     // 9h
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimeSpentColor(tt.seconds, hours)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetWorklogsFromWorklogIssues(t *testing.T) {
	worklogIssues := []*WorklogIssue{
		{
			Worklog: &Worklog{TimeSpentSeconds: 3600},
			Issue:   Issue{Key: "TEST-1"},
		},
		{
			Worklog: &Worklog{TimeSpentSeconds: 1800},
			Issue:   Issue{Key: "TEST-2"},
		},
	}

	result := getWorklogsFromWorklogIssues(worklogIssues)
	assert.Len(t, result, 2)
	assert.Equal(t, 3600, result[0].TimeSpentSeconds)
	assert.Equal(t, 1800, result[1].TimeSpentSeconds)
}

func TestMonthRange(t *testing.T) { //nolint:dupl
	testTime := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)

	start, end := MonthRange(&testTime)

	assert.Equal(t, 2024, start.Year())
	assert.Equal(t, time.Month(3), start.Month())
	assert.Equal(t, 1, start.Day())
	assert.Equal(t, 0, start.Hour())
	assert.Equal(t, 0, start.Minute())
	assert.Equal(t, 0, start.Second())

	assert.Equal(t, 2024, end.Year())
	assert.Equal(t, time.Month(3), end.Month())
	assert.Equal(t, 31, end.Day())
	assert.Equal(t, 23, end.Hour())
	assert.Equal(t, 59, end.Minute())
	assert.Equal(t, 59, end.Second())
}

func TestDayRange(t *testing.T) { //nolint:dupl
	testTime := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)

	start, end := DayRange(&testTime)

	assert.Equal(t, 2024, start.Year())
	assert.Equal(t, time.Month(3), start.Month())
	assert.Equal(t, 15, start.Day())
	assert.Equal(t, 0, start.Hour())
	assert.Equal(t, 0, start.Minute())
	assert.Equal(t, 0, start.Second())

	assert.Equal(t, 2024, end.Year())
	assert.Equal(t, time.Month(3), end.Month())
	assert.Equal(t, 15, end.Day())
	assert.Equal(t, 23, end.Hour())
	assert.Equal(t, 59, end.Minute())
	assert.Equal(t, 59, end.Second())
}

func TestFindIssueKeyInString_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no issue key", "just some random text", ""},
		{"lowercase", "ticket-123", ""}, // should not match lowercase
		{"partial match", "ABC", ""},    // needs number part
		{"number only", "123", ""},      // needs letter part
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindIssueKeyInString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
