package gojira

import (
	"bytes"
	"encoding/json"
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
	WorkLog WorkLog
	Issue   Issue
}

func GetWorkLogs() []WorkLog {
	currentDay := time.Now().Format("2006-01-02")
	requestUrl := fmt.Sprintf("%s/worklogs/user/%s?from=%s&to=%s", Config.TempoUrl, Config.JiraAccountId, currentDay, currentDay)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
		"Content-Type":  "application/json",
	}
	response := SendHttpRequest("GET", requestUrl, nil, headers, 200)
	var workLogsResponse WorkLogsResponse
	err := json.Unmarshal(response, &workLogsResponse)
	if err != nil {
		panic(err)
	}
	return workLogsResponse.WorkLogs
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

func (workLog *WorkLog) Update(timeSpent string) {
	workLog.TimeSpentSeconds = TimeSpentToSeconds(timeSpent)
	payload := WorkLogUpdate{
		IssueKey:         workLog.Issue.Key,
		StartDate:        workLog.StartDate,
		StartTime:        workLog.StartTime,
		Description:      workLog.Description,
		AuthorAccountId:  workLog.Author.AccountId,
		TimeSpentSeconds: workLog.TimeSpentSeconds,
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, workLog.TempoWorklogid)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
		"Content-Type":  "application/json",
	}

	SendHttpRequest("PUT", requestUrl, requestBody, headers, 200)
}
