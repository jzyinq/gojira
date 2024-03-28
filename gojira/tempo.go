package gojira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type TempoClient struct {
	Url           string
	Token         string
	JiraAccountId string
}

func NewTempoClient() *TempoClient {
	return &TempoClient{
		Url:           Config.TempoUrl,
		Token:         Config.TempoToken,
		JiraAccountId: Config.JiraAccountId,
	}
}

type WorklogsResponse struct {
	Worklogs []Worklog `json:"results"`
}

type WorklogUpdateRequest struct {
	IssueKey         string `json:"issueKey"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	AuthorAccountId  string `json:"authorAccountId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func (tc *TempoClient) GetWorklogs(fromDate, toDate time.Time) (WorklogsResponse, error) {
	// tempo is required only because of fetching worklogs by date range
	requestUrl := fmt.Sprintf("%s/worklogs/user/%s?from=%s&to=%s&limit=1000",
		tc.Url, tc.JiraAccountId, fromDate.Format(dateLayout), toDate.Format(dateLayout))
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", tc.Token),
		"Content-Type":  "application/json",
	}
	response, err := SendHttpRequest("GET", requestUrl, nil, headers, 200)
	if err != nil {
		return WorklogsResponse{}, err
	}
	var workLogsResponse WorklogsResponse
	err = json.Unmarshal(response, &workLogsResponse)
	if err != nil {
		return WorklogsResponse{}, err
	}
	return workLogsResponse, err
}

func (tc *TempoClient) UpdateWorklog(worklog *Worklog, timeSpent string) error {
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)

	payload := WorklogUpdateRequest{
		IssueKey:         worklog.Issue.Key,
		StartDate:        worklog.StartDate,
		StartTime:        worklog.StartTime,
		Description:      worklog.Description,
		AuthorAccountId:  worklog.Author.AccountId,
		TimeSpentSeconds: timeSpentInSeconds,
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, worklog.TempoWorklogid)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
		"Content-Type":  "application/json",
	}
	_, err := SendHttpRequest("PUT", requestUrl, requestBody, headers, 200)
	return err
}

func (tc *TempoClient) DeleteWorklog(tempoWorklogID int) error {
	requestUrl := fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, tempoWorklogID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", Config.TempoToken),
		"Content-Type":  "application/json",
	}
	_, err := SendHttpRequest("DELETE", requestUrl, nil, headers, 204)
	return err
}
