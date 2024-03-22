package gojira

import (
	"errors"
	"github.com/sirupsen/logrus"
	"log"
	"regexp"
	"strconv"
	"time"
)

func NewWorkLog(issueKey string, logTime *time.Time, timeSpent string) (WorkLog, error) {
	workLogResponse, err := NewJiraClient().CreateWorklog(issueKey, logTime, timeSpent)
	if err != nil {
		return WorkLog{}, err
	}

	jiraWorklogId, err := strconv.Atoi(workLogResponse.ID)
	if err != nil {
		return WorkLog{}, err
	}

	worklog := WorkLog{
		JiraWorklogid: jiraWorklogId,
		StartDate:     logTime.Format(dateLayout),
		StartTime:     logTime.Format("15:04:05"),
		Author: struct { // FIXME oh my god what a mess
			Self        string `json:"self"`
			AccountId   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		}{Self: workLogResponse.Self, AccountId: workLogResponse.Author.Accountid, DisplayName: workLogResponse.Author.Displayname},
		TimeSpentSeconds: workLogResponse.Timespentseconds,
		Issue: struct { // FIXME oh my god what a mess
			Self string `json:"self"`
			Key  string `json:"key"`
			ID   int    `json:"id"`
		}{Self: "", Key: issueKey, ID: 0},
	}
	return worklog, nil
}

type WorkLog struct {
	Self           string `json:"self"`
	TempoWorklogid int    `json:"tempoWorklogId"`
	JiraWorklogid  int    `json:"jiraWorklogId"`
	Issue          struct {
		Self string `json:"self"`
		Key  string `json:"key"`
		ID   int    `json:"id"`
	} `json:"issue"`
	TimeSpentSeconds int       `json:"timeSpentSeconds"`
	BillableSeconds  int       `json:"billableSeconds"`
	StartDate        string    `json:"startDate"`
	StartTime        string    `json:"startTime"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Author           struct {
		Self        string `json:"self"`
		AccountId   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	} `json:"author"`
	Attributes struct {
		Self   string        `json:"self"`
		Values []interface{} `json:"values"`
	} `json:"attributes"`
}

type WorkLogIssue struct {
	WorkLog *WorkLog
	Issue   Issue
}

type WorkLogsIssues struct {
	startDate time.Time
	endDate   time.Time
	issues    []WorkLogIssue
}

type WorkLogs struct {
	startDate time.Time
	endDate   time.Time
	logs      []*WorkLog
}

func (wl *WorkLogs) LogsOnDate(date *time.Time) ([]*WorkLog, error) {
	var logsOnDate []*WorkLog
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

// function that will summarize all timeSpentSeconds in logs slice and return it
func (wl *WorkLogs) TotalTimeSpent() int {
	var totalTime int
	for _, log := range wl.logs {
		totalTime += log.TimeSpentSeconds
	}
	return totalTime
}

func (wli *WorkLogsIssues) IssuesOnDate(date *time.Time) ([]*WorkLogIssue, error) {
	var issuesOnDate []*WorkLogIssue
	if date.Before(wli.startDate) || date.After(wli.endDate) {
		return nil, errors.New("Date is out of worklogs range")
	}
	truncatedDate := (*date).Truncate(24 * time.Hour)
	for i, issue := range wli.issues {
		// FIXME should be in local timezone PariseInLocation - but it's not working
		logDate, err := time.Parse(dateLayout, issue.WorkLog.StartDate)
		logDate = logDate.Truncate(24 * time.Hour)
		if err != nil {
			return nil, err
		}
		if truncatedDate.Equal(logDate.Truncate(24 * time.Hour)) {
			logrus.Info("truncatedDate ", truncatedDate)
			issuesOnDate = append(issuesOnDate, &wli.issues[i])
		}
	}
	return issuesOnDate, nil
}

func GetWorkLogs() (WorkLogs, error) {
	// get first day of week nd the last for date in app.time
	fromDate, toDate := MonthRange(app.time)
	logrus.Debug("getting worklogs from %s to %s...", fromDate, toDate)
	workLogsResponse, err := NewTempoClient().GetWorklogs(fromDate, toDate)
	if err != nil {
		return WorkLogs{}, err
	}
	var worklogs []*WorkLog
	for i, _ := range workLogsResponse.WorkLogs {
		worklogs = append(worklogs, &workLogsResponse.WorkLogs[i])
	}
	return WorkLogs{startDate: fromDate, endDate: toDate, logs: worklogs}, nil
}

func TimeSpentToSeconds(timeSpent string) int {
	r, _ := regexp.Compile("(([0-9]+)h)?\\s?(([0-9]+)m)?")
	match := r.FindStringSubmatch(timeSpent)
	var timeSpentSeconds int = 0

	if match[1] != "" {
		hours, err := strconv.ParseInt(match[2], 10, 64)
		timeSpentSeconds += int(hours) * 60 * 60
		if err != nil {
			log.Fatal(err)
		}
	}
	if match[3] != "" {
		minutes, err := strconv.ParseInt(match[4], 10, 32)
		timeSpentSeconds += int(minutes) * 60
		if err != nil {
			log.Fatal(err)
		}
	}
	return timeSpentSeconds
}

func (wl *WorkLog) Update(timeSpent string) error {
	logrus.Debugf("updating worklog ... %+v", wl)
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)
	var err error

	if wl.TempoWorklogid != 0 {
		// make update request to tempo if tempoWorklogId is set
		err = NewTempoClient().UpdateWorklog(wl, timeSpent)
	} else {
		// make update request to jira if tempoWorklogId is not set
		err = NewJiraClient().UpdateWorklog(wl.Issue.Key, wl.JiraWorklogid, timeSpentInSeconds)
	}
	if err != nil {
		return err
	}
	wl.TimeSpentSeconds = timeSpentInSeconds
	return nil
}

func (wl *WorkLogs) Delete(worklog *WorkLog) error {
	logrus.Debugf("deleting worklog ... %+v", worklog)
	// make update request to tempo if tempoWorklogId is set
	var err error
	if worklog.TempoWorklogid != 0 {
		err = NewTempoClient().DeleteWorklog(worklog.TempoWorklogid)
	} else {
		err = NewJiraClient().DeleteWorklog(worklog.Issue.Key, worklog.JiraWorklogid)
	}
	if err != nil {
		logrus.Debug(worklog)
		return err
	}

	// FIXME delete is kinda buggy - it messes up pointers and we're getting weird results
	for i, issue := range app.workLogsIssues.issues {
		if issue.WorkLog.JiraWorklogid == worklog.JiraWorklogid {
			app.workLogsIssues.issues = append(app.workLogsIssues.issues[:i], app.workLogsIssues.issues[i+1:]...)
			break
		}
	}
	for i, log := range wl.logs {
		if log.JiraWorklogid == worklog.JiraWorklogid {
			wl.logs = append(wl.logs[:i], wl.logs[i+1:]...)
			break
		}
	}

	return nil
}
