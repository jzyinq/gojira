package gojira

import (
	"errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"time"
)

var timeSpentRegex = regexp.MustCompile(`(([0-9]+)h)?\s?(([0-9]+)m)?`)

func NewWorklog(issueId int, logTime *time.Time, timeSpent string) (Worklog, error) {
	workLogResponse, err := NewJiraClient().CreateWorklog(issueId, logTime, timeSpent)
	if err != nil {
		return Worklog{}, err
	}

	jiraWorklogID, err := strconv.Atoi(workLogResponse.ID)
	if err != nil {
		return Worklog{}, err
	}

	workLog := Worklog{
		JiraWorklogID:    jiraWorklogID,
		StartDate:        logTime.Format(dateLayout),
		StartTime:        logTime.Format("15:04:05"),
		TimeSpentSeconds: workLogResponse.Timespentseconds,
		Issue: struct {
			Id int `json:"id"`
		}{Id: issueId},
	}
	return workLog, nil
}

type Worklog struct {
	TempoWorklogid int `json:"tempoWorklogId"`
	JiraWorklogID  int `json:"jiraWorklogId"`
	Issue          struct {
		Id int `json:"id"`
	} `json:"issue"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	Author           struct {
		AccountId string `json:"accountId"`
	} `json:"author"`
}

type WorklogIssue struct {
	Worklog *Worklog
	Issue   Issue
}

type WorklogsIssues struct {
	startDate time.Time
	endDate   time.Time
	issues    []WorklogIssue
}

type Worklogs struct {
	startDate time.Time
	endDate   time.Time
	logs      []*Worklog
}

func (wl *Worklogs) LogsOnDate(date *time.Time) ([]*Worklog, error) {
	var logsOnDate []*Worklog
	truncatedDate := (*date).Truncate(24 * time.Hour)
	if truncatedDate.Before(wl.startDate) || truncatedDate.After(wl.endDate) {
		return nil, nil
	}
	for i, logEntry := range wl.logs {
		logDate, err := time.Parse(dateLayout, logEntry.StartDate)
		logDate = logDate.Truncate(24 * time.Hour)
		if err != nil {
			return nil, err
		}
		if truncatedDate.Equal(logDate) {
			logsOnDate = append(logsOnDate, wl.logs[i])
		}
	}
	return logsOnDate, nil
}

func findWorklogByIssueKey(worklogs []*Worklog, issueKey string) *Worklog {
	for _, log := range worklogs {
		if strconv.Itoa(log.Issue.Id) == issueKey {
			return log
		}
	}
	return nil
}

func GetIssuesWithWorklogs(worklogs []*Worklog) ([]Issue, error) {
	var err error
	var worklogIssueIds []int
	for _, worklog := range worklogs {
		worklogIssueIds = append(worklogIssueIds, worklog.Issue.Id)
	}
	if len(worklogIssueIds) == 0 {
		return []Issue{}, err
	}
	todaysIssues, err := NewJiraClient().GetIssuesByKeys(worklogIssueIds)
	if err != nil {
		return []Issue{}, err
	}
	return todaysIssues.Issues, nil
}

func (wl *Worklogs) TotalTimeSpentToPresentDay() int {
	totalTime := 0
	for _, log := range wl.logs {
		logDate, err := time.Parse(dateLayout, log.StartDate)
		if err != nil {
			logrus.Error(err)
		}
		if logDate.Before(time.Now().UTC()) {
			totalTime += log.TimeSpentSeconds
		}
	}
	return totalTime
}

func (wli *WorklogsIssues) IssuesOnDate(date *time.Time) ([]*WorklogIssue, error) {
	var issuesOnDate []*WorklogIssue
	if date.Before(wli.startDate) || date.After(wli.endDate) {
		return nil, errors.New("Date is out of worklogs range")
	}
	truncatedDate := (*date).Truncate(24 * time.Hour)
	for i, issue := range wli.issues {
		logDate, err := time.Parse(dateLayout, issue.Worklog.StartDate)
		logDate = logDate.Truncate(24 * time.Hour)
		if err != nil {
			return nil, err
		}
		if truncatedDate.Equal(logDate.Truncate(24 * time.Hour)) {
			logrus.Debug("truncatedDate ", truncatedDate)
			issuesOnDate = append(issuesOnDate, &wli.issues[i])
		}
	}
	return issuesOnDate, nil
}

func GetWorklogs(fromDate time.Time, toDate time.Time) (Worklogs, error) {
	logrus.Infof("getting worklogs from %s to %s...", fromDate, toDate)
	workLogsResponse, err := NewTempoClient().GetWorklogs(fromDate, toDate)
	if err != nil {
		return Worklogs{}, err
	}
	var worklogs []*Worklog
	for i := range workLogsResponse.Worklogs {
		worklogs = append(worklogs, &workLogsResponse.Worklogs[i])
	}
	return Worklogs{startDate: fromDate, endDate: toDate, logs: worklogs}, nil
}

func TimeSpentToSeconds(timeSpent string) int {
	match := timeSpentRegex.FindStringSubmatch(timeSpent)
	var timeSpentSeconds int = 0

	if match[1] != "" {
		hours, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			logrus.Errorf("failed to parse hours from time spent '%s': %v", timeSpent, err)
			return 0
		}
		timeSpentSeconds += int(hours) * 60 * 60
	}
	if match[3] != "" {
		minutes, err := strconv.ParseInt(match[4], 10, 32)
		if err != nil {
			logrus.Errorf("failed to parse minutes from time spent '%s': %v", timeSpent, err)
			return 0
		}
		timeSpentSeconds += int(minutes) * 60
	}
	return timeSpentSeconds
}

func (wl *Worklog) Update(timeSpent string) error {
	logrus.Debugf("updating worklog ... %+v", wl)
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)
	var err error

	if wl.TempoWorklogid != 0 {
		// make update request to tempo if tempoWorklogId is set
		err = NewTempoClient().UpdateWorklog(wl, timeSpent)
	} else {
		// make update request to jira if tempoWorklogId is not set
		err = NewJiraClient().UpdateWorklog(wl.Issue.Id, wl.JiraWorklogID, timeSpentInSeconds)
	}
	if err != nil {
		return err
	}
	wl.TimeSpentSeconds = timeSpentInSeconds
	return nil
}

func (wl *Worklogs) Delete(w *Worklog) error {
	logrus.Debugf("deleting w ... %+v", w)
	// make update request to tempo if tempoWorklogId is set
	var err error
	if w.TempoWorklogid != 0 {
		err = NewTempoClient().DeleteWorklog(w.TempoWorklogid)
	} else {
		err = NewJiraClient().DeleteWorklog(w.Issue.Id, w.JiraWorklogID)
	}
	if err != nil {
		logrus.Debug(w)
		return err
	}

	// Remove from workLogsIssues by filtering
	app.mu.Lock()
	filtered := make([]WorklogIssue, 0, len(app.workLogsIssues.issues))
	for _, issue := range app.workLogsIssues.issues {
		if issue.Worklog.JiraWorklogID != w.JiraWorklogID {
			filtered = append(filtered, issue)
		}
	}
	app.workLogsIssues.issues = filtered

	// Remove from worklogs by filtering
	filteredLogs := make([]*Worklog, 0, len(wl.logs))
	for _, workLog := range wl.logs {
		if workLog.JiraWorklogID != w.JiraWorklogID {
			filteredLogs = append(filteredLogs, workLog)
		}
	}
	wl.logs = filteredLogs
	app.mu.Unlock()

	return nil
}
