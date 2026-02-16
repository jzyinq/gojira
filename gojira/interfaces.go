package gojira

import "time"

// JiraClientInterface defines the contract for interacting with Jira API
type JiraClientInterface interface {
	GetIssuesByJQL(jql string, maxResults int) (JQLResponse, error)
	GetLatestIssues() (JQLResponse, error)
	GetIssuesByKeys(issueKeys []int) (JQLResponse, error)
	GetIssue(issueKey string) (Issue, error)
	CreateWorklog(issueId int, logTime *time.Time, timeSpent string) (WorklogResponse, error)
	UpdateWorklog(issueId int, jiraWorklogId int, timeSpentInSeconds int) error
	DeleteWorklog(issueId int, jiraWorklogId int) error
}

// TempoClientInterface defines the contract for interacting with Tempo API
type TempoClientInterface interface {
	GetWorklogs(fromDate, toDate time.Time) (WorklogsResponse, error)
	UpdateWorklog(worklog *Worklog, timeSpent string) error
	DeleteWorklog(tempoWorklogID int) error
}

// Compile-time verification that concrete types implement interfaces
var _ JiraClientInterface = (*JiraClient)(nil)
var _ TempoClientInterface = (*TempoClient)(nil)
