package gojira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

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

type WorkLogUpdate struct {
	IssueKey         string `json:"issueKey"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	AuthorAccountId  string `json:"authorAccountId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

type JiraWorklogUpdate struct {
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}

type WorkLogsResponse struct {
	Self     string `json:"self"`
	Metadata struct {
		Count  int `json:"count"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"metadata"`
	WorkLogs []WorkLog `json:"results"`
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

func (w *WorkLogs) LogsOnDate(date time.Time) ([]*WorkLog, error) {
	var logsOnDate []*WorkLog
	if date.Before(w.startDate) || date.After(w.endDate) {
		return nil, nil
	}
	date = date.Truncate(24 * time.Hour)
	for i, log := range w.logs {
		logDate, err := time.ParseInLocation(dateLayout, log.StartDate, time.Local)
		logDate = logDate.Truncate(24 * time.Hour)
		if err != nil {
			return nil, err
		}
		if date.Equal(logDate.Truncate(24 * time.Hour)) {
			logsOnDate = append(logsOnDate, w.logs[i])
		}
	}
	return logsOnDate, nil
}

// function that will summarize all timeSpentSeconds in logs slice and return it
func (w *WorkLogs) TotalTimeSpent() int {
	var totalTime int
	for _, log := range w.logs {
		totalTime += log.TimeSpentSeconds
	}
	return totalTime
}

func (w *WorkLogsIssues) IssuesOnDate(date time.Time) ([]*WorkLogIssue, error) {
	var issuesOnDate []*WorkLogIssue
	if date.Before(w.startDate) || date.After(w.endDate) {
		return nil, errors.New("Date is out of worklogs range")
	}
	date = date.Truncate(24 * time.Hour)
	for i, issue := range w.issues {
		// FIXME should be in local timezone PariseInLocation - but it's not working
		logDate, err := time.Parse(dateLayout, issue.WorkLog.StartDate)
		logDate = logDate.Truncate(24 * time.Hour)
		if err != nil {
			return nil, err
		}
		if date.Equal(logDate.Truncate(24 * time.Hour)) {
			issuesOnDate = append(issuesOnDate, &w.issues[i])
		}
	}
	return issuesOnDate, nil
}

func GetWorkLogs() (WorkLogs, error) {
	// get first day of week nd the last for date in app.time
	fromDate, toDate := MonthRange(app.time)
	requestUrl := fmt.Sprintf("%s/worklogs/user/%s?from=%s&to=%s&limit=1000", Config.TempoUrl, Config.JiraAccountId, fromDate.Format(dateLayout), toDate.Format(dateLayout))
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
		"Content-Type":  "application/json",
	}
	response, err := SendHttpRequest("GET", requestUrl, nil, headers, 200)
	if err != nil {
		return WorkLogs{}, err
	}
	var workLogsResponse WorkLogsResponse
	err = json.Unmarshal(response, &workLogsResponse)
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

func (workLog *WorkLog) Update(timeSpent string) error {
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)

	payload := JiraWorklogUpdate{
		TimeSpentSeconds: timeSpentInSeconds,
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false", Config.JiraUrl, workLog.Issue.Key, workLog.JiraWorklogid)
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}

	_, err := SendHttpRequest("PUT", requestUrl, requestBody, headers, 200)
	if err != nil {
		return err
	}

	workLog.TimeSpentSeconds = timeSpentInSeconds
	return nil
}

func (workLogs *WorkLogs) Delete(worklog *WorkLog) error {
	// fixme this no longer depends on tempo api - why I'm using tempo api? :confused:
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false",
		Config.JiraUrl, worklog.Issue.Key, worklog.JiraWorklogid)
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}

	_, err := SendHttpRequest("DELETE", requestUrl, nil, headers, 204)
	if err != nil {
		return err
	}

	// FIXME delete is kinda buggy - it messes up pointers and we're getting weird results\
	for i, issue := range app.workLogsIssues.issues {
		if issue.WorkLog.JiraWorklogid == worklog.JiraWorklogid {
			app.workLogsIssues.issues = append(app.workLogsIssues.issues[:i], app.workLogsIssues.issues[i+1:]...)
			break
		}
	}
	for i, log := range workLogs.logs {
		if log.JiraWorklogid == worklog.JiraWorklogid {
			workLogs.logs = append(workLogs.logs[:i], workLogs.logs[i+1:]...)
			break
		}
	}

	return nil
}
