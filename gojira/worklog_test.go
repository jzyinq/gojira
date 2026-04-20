package gojira

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorklog_TimeSpentToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"hours and minutes", "1h 30m", 5400},
		{"hours and minutes no space", "1h30m", 5400},
		{"only minutes", "29m", 1740},
		{"only hours", "2h", 7200},
		{"multiple hours and minutes", "3h 45m", 13500},
		{"empty string", "", 0},
		{"just spaces", "   ", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeSpentToSeconds(tt.input)
			assert.Equal(t, tt.expected, result, "TimeSpentToSeconds(%s) should return %d", tt.input, tt.expected)
		})
	}
}

func TestWorklog_LogsOnDate(t *testing.T) {
	// Setup test data
	date1, _ := time.Parse(dateLayout, "2024-01-15")
	date2, _ := time.Parse(dateLayout, "2024-01-16")
	date3, _ := time.Parse(dateLayout, "2024-01-17")

	worklogs := Worklogs{
		startDate: date1,
		endDate:   date3,
		logs: []*Worklog{
			{StartDate: "2024-01-15", TimeSpentSeconds: 3600},
			{StartDate: "2024-01-15", TimeSpentSeconds: 1800},
			{StartDate: "2024-01-16", TimeSpentSeconds: 7200},
			{StartDate: "2024-01-17", TimeSpentSeconds: 900},
		},
	}

	t.Run("returns logs for date with multiple entries", func(t *testing.T) {
		logs, err := worklogs.LogsOnDate(&date1)
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		assert.Equal(t, 3600, logs[0].TimeSpentSeconds)
		assert.Equal(t, 1800, logs[1].TimeSpentSeconds)
	})

	t.Run("returns logs for date with single entry", func(t *testing.T) {
		logs, err := worklogs.LogsOnDate(&date2)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, 7200, logs[0].TimeSpentSeconds)
	})

	t.Run("returns empty for date before start", func(t *testing.T) {
		beforeDate, _ := time.Parse(dateLayout, "2024-01-14")
		logs, err := worklogs.LogsOnDate(&beforeDate)
		assert.NoError(t, err)
		assert.Nil(t, logs)
	})

	t.Run("returns empty for date after end", func(t *testing.T) {
		afterDate, _ := time.Parse(dateLayout, "2024-01-18")
		logs, err := worklogs.LogsOnDate(&afterDate)
		assert.NoError(t, err)
		assert.Nil(t, logs)
	})
}

func TestWorklog_TotalTimeSpentToPresentDay(t *testing.T) {
	// Create dates relative to now
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)

	worklogs := Worklogs{
		startDate: yesterday,
		endDate:   tomorrow,
		logs: []*Worklog{
			{StartDate: yesterday.Format(dateLayout), TimeSpentSeconds: 3600}, // 1 hour
			{StartDate: yesterday.Format(dateLayout), TimeSpentSeconds: 1800}, // 30 minutes
			{StartDate: tomorrow.Format(dateLayout), TimeSpentSeconds: 7200},  // 2 hours (future, should be excluded)
		},
	}

	total := worklogs.TotalTimeSpentToPresentDay()

	// Should only count past entries (yesterday's entries = 5400 seconds)
	assert.Equal(t, 5400, total, "Should only count worklogs before present")
}

func TestWorklogsIssues_IssuesOnDate(t *testing.T) {
	date1, _ := time.Parse(dateLayout, "2024-01-15")
	date3, _ := time.Parse(dateLayout, "2024-01-17")

	worklogsIssues := WorklogsIssues{
		startDate: date1,
		endDate:   date3,
		issues: []WorklogIssue{
			{
				Worklog: &Worklog{StartDate: "2024-01-15", TimeSpentSeconds: 3600},
				Issue:   Issue{Key: "TEST-1"},
			},
			{
				Worklog: &Worklog{StartDate: "2024-01-15", TimeSpentSeconds: 1800},
				Issue:   Issue{Key: "TEST-2"},
			},
			{
				Worklog: &Worklog{StartDate: "2024-01-16", TimeSpentSeconds: 7200},
				Issue:   Issue{Key: "TEST-3"},
			},
		},
	}

	t.Run("returns issues for valid date", func(t *testing.T) {
		issues, err := worklogsIssues.IssuesOnDate(&date1)
		assert.NoError(t, err)
		assert.Len(t, issues, 2)
		assert.Equal(t, "TEST-1", issues[0].Issue.Key)
		assert.Equal(t, "TEST-2", issues[1].Issue.Key)
	})

	t.Run("returns error for date before range", func(t *testing.T) {
		beforeDate := date1.AddDate(0, 0, -1)
		_, err := worklogsIssues.IssuesOnDate(&beforeDate)
		assert.Error(t, err)
		assert.Equal(t, "date is out of worklogs range", err.Error())
	})

	t.Run("returns error for date after range", func(t *testing.T) {
		afterDate := date3.AddDate(0, 0, 1)
		_, err := worklogsIssues.IssuesOnDate(&afterDate)
		assert.Error(t, err)
		assert.Equal(t, "date is out of worklogs range", err.Error())
	})
}

func TestFindWorklogByIssueKey(t *testing.T) {
	worklogs := []*Worklog{
		{Issue: struct {
			Id int `json:"id"`
		}{Id: 123}, TimeSpentSeconds: 3600},
		{Issue: struct {
			Id int `json:"id"`
		}{Id: 456}, TimeSpentSeconds: 1800},
		{Issue: struct {
			Id int `json:"id"`
		}{Id: 789}, TimeSpentSeconds: 7200},
	}

	t.Run("finds worklog by matching issue id", func(t *testing.T) {
		result := findWorklogByIssueKey(worklogs, "456")
		assert.NotNil(t, result)
		assert.Equal(t, 456, result.Issue.Id)
		assert.Equal(t, 1800, result.TimeSpentSeconds)
	})

	t.Run("returns nil for non-existent issue", func(t *testing.T) {
		result := findWorklogByIssueKey(worklogs, "999")
		assert.Nil(t, result)
	})

	t.Run("returns nil for empty worklogs", func(t *testing.T) {
		result := findWorklogByIssueKey([]*Worklog{}, "123")
		assert.Nil(t, result)
	})
}

func TestIssue_GetIdAsInt(t *testing.T) {
	t.Run("parses valid numeric id", func(t *testing.T) {
		issue := Issue{Id: "12345"}
		assert.Equal(t, 12345, issue.GetIdAsInt())
	})

	t.Run("returns 0 for invalid id", func(t *testing.T) {
		issue := Issue{Id: "not-a-number"}
		assert.Equal(t, 0, issue.GetIdAsInt())
	})

	t.Run("returns 0 for empty id", func(t *testing.T) {
		issue := Issue{Id: ""}
		assert.Equal(t, 0, issue.GetIdAsInt())
	})
}

func TestRemoveWorklog(t *testing.T) { //nolint:funlen
	logs := []*Worklog{
		{JiraWorklogID: 1, TimeSpentSeconds: 3600},
		{JiraWorklogID: 2, TimeSpentSeconds: 1800},
		{JiraWorklogID: 3, TimeSpentSeconds: 900},
	}
	issues := []WorklogIssue{
		{Worklog: logs[0], Issue: Issue{Key: "TEST-1"}},
		{Worklog: logs[1], Issue: Issue{Key: "TEST-2"}},
		{Worklog: logs[2], Issue: Issue{Key: "TEST-3"}},
	}

	t.Run("removes matching worklog from both slices", func(t *testing.T) {
		filteredLogs, filteredIssues := removeWorklog(logs, issues, logs[1])
		assert.Len(t, filteredLogs, 2)
		assert.Len(t, filteredIssues, 2)
		assert.Equal(t, 1, filteredLogs[0].JiraWorklogID)
		assert.Equal(t, 3, filteredLogs[1].JiraWorklogID)
		assert.Equal(t, "TEST-1", filteredIssues[0].Issue.Key)
		assert.Equal(t, "TEST-3", filteredIssues[1].Issue.Key)
	})

	t.Run("returns same contents when pointer not found", func(t *testing.T) {
		other := &Worklog{JiraWorklogID: 999}
		filteredLogs, filteredIssues := removeWorklog(logs, issues, other)
		assert.Len(t, filteredLogs, 3)
		assert.Len(t, filteredIssues, 3)
	})

	t.Run("handles empty slices", func(t *testing.T) {
		filteredLogs, filteredIssues := removeWorklog([]*Worklog{}, []WorklogIssue{}, logs[0])
		assert.Empty(t, filteredLogs)
		assert.Empty(t, filteredIssues)
	})

	t.Run("removes first element", func(t *testing.T) {
		filteredLogs, filteredIssues := removeWorklog(logs, issues, logs[0])
		assert.Len(t, filteredLogs, 2)
		assert.Len(t, filteredIssues, 2)
		assert.Equal(t, 2, filteredLogs[0].JiraWorklogID)
		assert.Equal(t, 3, filteredLogs[1].JiraWorklogID)
	})

	t.Run("removes last element", func(t *testing.T) {
		filteredLogs, filteredIssues := removeWorklog(logs, issues, logs[2])
		assert.Len(t, filteredLogs, 2)
		assert.Len(t, filteredIssues, 2)
		assert.Equal(t, 1, filteredLogs[0].JiraWorklogID)
		assert.Equal(t, 2, filteredLogs[1].JiraWorklogID)
	})

	t.Run("does not mutate original slices", func(t *testing.T) {
		originalLen := len(logs)
		removeWorklog(logs, issues, logs[1])
		assert.Len(t, logs, originalLen, "original logs slice should not be modified")
	})

	t.Run("removes only the exact pointer when IDs are zero", func(t *testing.T) {
		zeroLogs := []*Worklog{{JiraWorklogID: 0}, {JiraWorklogID: 0}, {JiraWorklogID: 0}}
		zeroIssues := []WorklogIssue{
			{Worklog: zeroLogs[0]},
			{Worklog: zeroLogs[1]},
			{Worklog: zeroLogs[2]},
		}
		filteredLogs, filteredIssues := removeWorklog(zeroLogs, zeroIssues, zeroLogs[1])
		assert.Len(t, filteredLogs, 2)
		assert.Len(t, filteredIssues, 2)
		assert.Same(t, zeroLogs[0], filteredLogs[0])
		assert.Same(t, zeroLogs[2], filteredLogs[1])
	})
}
