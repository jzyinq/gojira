package gojira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveLogArgs(t *testing.T) {
	tests := []struct {
		name           string
		arg0           string
		arg1           string
		gitIssueKey    string
		wantIssueKey   string
		wantTimeSpent  string
	}{
		{
			name:          "issue and time provided explicitly",
			arg0:          "TICKET-123",
			arg1:          "30m",
			gitIssueKey:   "GIT-999",
			wantIssueKey:  "TICKET-123",
			wantTimeSpent: "30m",
		},
		{
			name:          "only time provided, issue from git branch",
			arg0:          "30m",
			arg1:          "",
			gitIssueKey:   "GIT-999",
			wantIssueKey:  "GIT-999",
			wantTimeSpent: "30m",
		},
		{
			name:          "issue URL with time provided explicitly",
			arg0:          "https://instance.atlassian.net/browse/TICKET-123",
			arg1:          "1h30m",
			gitIssueKey:   "GIT-999",
			wantIssueKey:  "TICKET-123",
			wantTimeSpent: "1h30m",
		},
		{
			name:          "only issue provided, no time",
			arg0:          "TICKET-123",
			arg1:          "",
			gitIssueKey:   "GIT-999",
			wantIssueKey:  "TICKET-123",
			wantTimeSpent: "",
		},
		{
			name:          "time provided but no issue in arg or git branch",
			arg0:          "30m",
			arg1:          "",
			gitIssueKey:   "",
			wantIssueKey:  "",
			wantTimeSpent: "30m",
		},
		{
			name:          "no args, no git branch",
			arg0:          "",
			arg1:          "",
			gitIssueKey:   "",
			wantIssueKey:  "",
			wantTimeSpent: "",
		},
		{
			name:          "no args, issue from git branch",
			arg0:          "",
			arg1:          "",
			gitIssueKey:   "GIT-999",
			wantIssueKey:  "GIT-999",
			wantTimeSpent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issueKey, timeSpent := resolveLogArgs(tt.arg0, tt.arg1, tt.gitIssueKey)
			assert.Equal(t, tt.wantIssueKey, issueKey)
			assert.Equal(t, tt.wantTimeSpent, timeSpent)
		})
	}
}
