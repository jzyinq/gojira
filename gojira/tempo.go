package gojira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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

type JiraWorklogUpdate struct {
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}

type WorkLogUpdate struct {
	IssueKey         string `json:"issueKey"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	AuthorAccountId  string `json:"authorAccountId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
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
			issuesOnDate = append(issuesOnDate, &wli.issues[i])
		}
	}
	return issuesOnDate, nil
}

func GetWorkLogs() (WorkLogs, error) {
	// get first day of week nd the last for date in app.time
	fromDate, toDate := MonthRange(app.time)
	// tempo is required only for fetching workklogs by date range
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

func (wl *WorkLog) Update(timeSpent string) error {
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)

	// make update request to tempo if tempoWorklogId is set
	var requestBody *bytes.Buffer
	var requestUrl string
	var headers map[string]string

	// updating meetings does not work even through tempo? w00t
	if wl.TempoWorklogid != 0 {
		payload := WorkLogUpdate{
			IssueKey:         wl.Issue.Key,
			StartDate:        wl.StartDate,
			StartTime:        wl.StartTime,
			Description:      wl.Description,
			AuthorAccountId:  wl.Author.AccountId,
			TimeSpentSeconds: timeSpentInSeconds,
		}
		payloadJson, _ := json.Marshal(payload)
		requestBody = bytes.NewBuffer(payloadJson)
		requestUrl = fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, wl.TempoWorklogid)
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
			"Content-Type":  "application/json",
		}
	} else {
		payload := JiraWorklogUpdate{
			TimeSpentSeconds: timeSpentInSeconds,
		}
		payloadJson, _ := json.Marshal(payload)
		requestBody = bytes.NewBuffer(payloadJson)
		// FIXME use tempo api to update worklog, unless there is not tempoId in worklog
		requestUrl = fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false", Config.JiraUrl, wl.Issue.Key, wl.JiraWorklogid)
		headers = map[string]string{
			"Authorization": getJiraAuthorizationHeader(),
			"Content-Type":  "application/json",
		}
	}
	_, err := SendHttpRequest("PUT", requestUrl, requestBody, headers, 200)
	if err != nil {
		return err
	}

	wl.TimeSpentSeconds = timeSpentInSeconds
	return nil
}

func (wl *WorkLogs) Delete(worklog *WorkLog) error {
	logrus.Infof("Deleting worklog ... %+v", worklog)
	// make update request to tempo if tempoWorklogId is set
	var requestUrl string
	var headers map[string]string

	if worklog.TempoWorklogid != 0 {
		requestUrl = fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, worklog.TempoWorklogid)
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
			"Content-Type":  "application/json",
		}
	} else {
		requestUrl = fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false", Config.JiraUrl, worklog.Issue.Key, worklog.JiraWorklogid)
		headers = map[string]string{
			"Authorization": getJiraAuthorizationHeader(),
			"Content-Type":  "application/json",
		}
	}
	_, err := SendHttpRequest("DELETE", requestUrl, nil, headers, 204)
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
