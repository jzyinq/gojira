package gojira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testIssueKey    = "TICKET-123"
	testDuration30m = "30m"
	testGitIssueKey = "GIT-999"
)

func TestResolveLogArgs(t *testing.T) { //nolint:funlen
	tests := []struct {
		name          string
		arg0          string
		arg1          string
		gitIssueKey   string
		wantIssueKey  string
		wantTimeSpent string
	}{
		{
			name:          "issue and time provided explicitly",
			arg0:          testIssueKey,
			arg1:          testDuration30m,
			gitIssueKey:   testGitIssueKey,
			wantIssueKey:  testIssueKey,
			wantTimeSpent: testDuration30m,
		},
		{
			name:          "only time provided, issue from git branch",
			arg0:          testDuration30m,
			arg1:          "",
			gitIssueKey:   testGitIssueKey,
			wantIssueKey:  testGitIssueKey,
			wantTimeSpent: testDuration30m,
		},
		{
			name:          "issue URL with time provided explicitly",
			arg0:          "https://instance.atlassian.net/browse/TICKET-123",
			arg1:          testDuration1h30m,
			gitIssueKey:   testGitIssueKey,
			wantIssueKey:  testIssueKey,
			wantTimeSpent: testDuration1h30m,
		},
		{
			name:          "only issue provided, no time",
			arg0:          testIssueKey,
			arg1:          "",
			gitIssueKey:   testGitIssueKey,
			wantIssueKey:  testIssueKey,
			wantTimeSpent: "",
		},
		{
			name:          "time provided but no issue in arg or git branch",
			arg0:          testDuration30m,
			arg1:          "",
			gitIssueKey:   "",
			wantIssueKey:  "",
			wantTimeSpent: testDuration30m,
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
			gitIssueKey:   testGitIssueKey,
			wantIssueKey:  testGitIssueKey,
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
