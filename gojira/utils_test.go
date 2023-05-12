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
	fixture := []WorkLog{
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
		"https://instance.atlassian.net/secure/RapidBoard.jspa?rapidView=84&projectKey=TICKET&view=planning&selectedIssue=TICKET-999&issueLimit=100",
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
		{"1h30m", 5400},
		{"29m", 1740},
	}

	for _, fixture := range fixtures {
		timeSpentInSeconds := TimeSpentToSeconds(fixture.TimeSpent)
		if timeSpentInSeconds != fixture.expectedTimeSpentInSeconds {
			t.Errorf("Incorrect timeSpent - got %d instead of %d", timeSpentInSeconds, fixture.expectedTimeSpentInSeconds)
		}
	}
}
