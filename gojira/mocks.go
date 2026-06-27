package gojira

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockJiraClient is a mock implementation of JiraClientInterface
type MockJiraClient struct {
	mock.Mock
}

func (m *MockJiraClient) GetIssuesByJQL(jql string, maxResults int) (JQLResponse, error) {
	args := m.Called(jql, maxResults)
	return args.Get(0).(JQLResponse), args.Error(1)
}

func (m *MockJiraClient) GetLatestIssues() (JQLResponse, error) {
	args := m.Called()
	return args.Get(0).(JQLResponse), args.Error(1)
}

func (m *MockJiraClient) GetIssuesByKeys(issueKeys []int) (JQLResponse, error) {
	args := m.Called(issueKeys)
	return args.Get(0).(JQLResponse), args.Error(1)
}

func (m *MockJiraClient) GetIssue(issueKey string) (Issue, error) {
	args := m.Called(issueKey)
	return args.Get(0).(Issue), args.Error(1)
}

func (m *MockJiraClient) CreateWorklog(issueId int, logTime *time.Time, timeSpent string) (WorklogResponse, error) {
	args := m.Called(issueId, logTime, timeSpent)
	return args.Get(0).(WorklogResponse), args.Error(1)
}

func (m *MockJiraClient) UpdateWorklog(issueId int, jiraWorklogId int, timeSpentInSeconds int) error {
	args := m.Called(issueId, jiraWorklogId, timeSpentInSeconds)
	return args.Error(0)
}

func (m *MockJiraClient) DeleteWorklog(issueId int, jiraWorklogId int) error {
	args := m.Called(issueId, jiraWorklogId)
	return args.Error(0)
}

// MockTempoClient is a mock implementation of TempoClientInterface
type MockTempoClient struct {
	mock.Mock
}

func (m *MockTempoClient) GetWorklogs(fromDate, toDate time.Time) (WorklogsResponse, error) {
	args := m.Called(fromDate, toDate)
	return args.Get(0).(WorklogsResponse), args.Error(1)
}

func (m *MockTempoClient) UpdateWorklog(worklog *Worklog, timeSpent string) error {
	args := m.Called(worklog, timeSpent)
	return args.Error(0)
}

func (m *MockTempoClient) DeleteWorklog(tempoWorklogID int) error {
	args := m.Called(tempoWorklogID)
	return args.Error(0)
}
