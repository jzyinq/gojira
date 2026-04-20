package gojira

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	Id               string `json:"id"`
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
	headers := CreateBearerAuthHeaders(tc.Token)
	response, err := SendHttpRequest("GET", requestUrl, nil, headers, 200)
	if err != nil {
		return WorklogsResponse{}, fmt.Errorf("failed to fetch worklogs from %s to %s: %w",
			fromDate.Format(dateLayout), toDate.Format(dateLayout), err)
	}
	var workLogsResponse WorklogsResponse
	err = json.Unmarshal(response, &workLogsResponse)
	if err != nil {
		return WorklogsResponse{}, fmt.Errorf("failed to unmarshal worklogs response: %w", err)
	}
	return workLogsResponse, err
}

func (tc *TempoClient) UpdateWorklog(worklog *Worklog, timeSpent string) error {
	timeSpentInSeconds := TimeSpentToSeconds(timeSpent)

	payload := WorklogUpdateRequest{
		Id:               strconv.Itoa(worklog.Issue.Id),
		StartDate:        worklog.StartDate,
		StartTime:        worklog.StartTime,
		Description:      worklog.Description,
		AuthorAccountId:  worklog.Author.AccountId,
		TimeSpentSeconds: timeSpentInSeconds,
	}
	requestUrl := fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, worklog.TempoWorklogid)
	headers := CreateBearerAuthHeaders(Config.TempoToken)
	_, err := SendJSONRequest("PUT", requestUrl, payload, headers, 200)
	if err != nil {
		return fmt.Errorf("failed to update worklog %d: %w", worklog.TempoWorklogid, err)
	}
	return nil
}

func (tc *TempoClient) DeleteWorklog(tempoWorklogID int) error {
	requestUrl := fmt.Sprintf("%s/worklogs/%d", Config.TempoUrl, tempoWorklogID)
	headers := CreateBearerAuthHeaders(Config.TempoToken)
	_, err := SendHttpRequest("DELETE", requestUrl, nil, headers, 204)
	if err != nil {
		return fmt.Errorf("failed to delete worklog %d: %w", tempoWorklogID, err)
	}
	return nil
}
