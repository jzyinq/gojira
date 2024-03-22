package gojira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type JiraClient struct {
	Url       string
	Login     string
	Token     string
	AccountID string
}

func NewJiraClient() *JiraClient {
	return &JiraClient{
		Url:       Config.JiraUrl,
		Login:     Config.JiraLogin,
		Token:     Config.JiraToken,
		AccountID: Config.JiraAccountId,
	}
}

func (jc *JiraClient) getHttpHeaders() map[string]string {

	authorizationToken := fmt.Sprintf("%s:%s", Config.JiraLogin, Config.JiraToken)
	authorizationHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authorizationToken)))
	return map[string]string{
		"Authorization": authorizationHeader,
		"Content-Type":  "application/json",
	}
}

type JQLSearch struct {
	Expand       []string `json:"expand"`
	Jql          string   `json:"jql"`
	MaxResults   int      `json:"maxResults"`
	FieldsByKeys bool     `json:"fieldsByKeys"`
	Fields       []string `json:"fields"`
	StartAt      int      `json:"startAt"`
}

type JQLResponse struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
	Names      struct {
		Summary string `json:"summary"`
		Status  string `json:"status"`
	} `json:"names"`
}

type Issue struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Self           string `json:"self"`
			Description    string `json:"description"`
			IconUrl        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				Self      string `json:"self"`
				ID        int    `json:"id"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
	} `json:"fields"`
}

func (jc *JiraClient) GetLatestIssues() (JQLResponse, error) {
	payload := &JQLSearch{
		Expand:       []string{"names"},
		Jql:          "assignee in (currentUser()) ORDER BY updated DESC, created DESC",
		MaxResults:   10,
		FieldsByKeys: false,
		Fields:       []string{"summary", "status"},
		StartAt:      0,
	}
	payloadJson, err := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/rest/api/2/search", Config.JiraUrl)
	response, err := SendHttpRequest("POST", requestUrl, requestBody, jc.getHttpHeaders(), 200)
	if err != nil {
		return JQLResponse{}, err
	}
	var jqlResponse JQLResponse
	err = json.Unmarshal(response, &jqlResponse)
	if err != nil {
		return JQLResponse{}, err
	}
	return jqlResponse, nil
}

func (jc *JiraClient) GetIssue(issueKey string) (Issue, error) {
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=summary,status", Config.JiraUrl, issueKey)
	response, err := SendHttpRequest("GET", requestUrl, nil, jc.getHttpHeaders(), 200)
	if err != nil {
		return Issue{}, err
	}
	var jiraIssue Issue
	err = json.Unmarshal(response, &jiraIssue)
	if err != nil {
		return Issue{}, err
	}
	return jiraIssue, nil
}

type WorkLogResponse struct {
	Self   string `json:"self"`
	Author struct {
		Self         string `json:"self"`
		Accountid    string `json:"accountId"`
		Emailaddress string `json:"emailAddress"`
		Avatarurls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		Displayname string `json:"displayName"`
		Active      bool   `json:"active"`
		Timezone    string `json:"timeZone"`
		Accounttype string `json:"accountType"`
	} `json:"author"`
	Updateauthor struct {
		Self         string `json:"self"`
		Accountid    string `json:"accountId"`
		Emailaddress string `json:"emailAddress"`
		Avatarurls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		Displayname string `json:"displayName"`
		Active      bool   `json:"active"`
		Timezone    string `json:"timeZone"`
		Accounttype string `json:"accountType"`
	} `json:"updateAuthor"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
	Started          string `json:"started"`
	Timespent        string `json:"timeSpent"`
	Timespentseconds int    `json:"timeSpentSeconds"`
	ID               string `json:"id"` // can it be an int? it's a number
	Issueid          string `json:"issueId"`
}

func (jc *JiraClient) CreateWorklog(issueKey string, logTime *time.Time, timeSpent string) (WorkLogResponse, error) {
	payload := map[string]string{
		"timeSpent":      timeSpent,
		"adjustEstimate": "leave",
		"started":        logTime.Format("2006-01-02T15:04:05.000-0700"),
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog?notifyUsers=false", Config.JiraUrl, issueKey)
	response, err := SendHttpRequest("POST", requestUrl, requestBody, jc.getHttpHeaders(), 201)
	if err != nil {
		return WorkLogResponse{}, err
	}

	var workLogRequest WorkLogResponse
	err = json.Unmarshal(response, &workLogRequest)
	if err != nil {
		return WorkLogResponse{}, err
	}
	return workLogRequest, nil
}

func (jc *JiraClient) UpdateWorklog(issueKey string, jiraWorklogId int, timeSpentInSeconds int) error {
	payload := JiraWorklogUpdate{
		TimeSpentSeconds: timeSpentInSeconds,
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	// FIXME use tempo api to update worklog, unless there is not tempoId in worklog
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false", Config.JiraUrl, issueKey, jiraWorklogId)
	_, err := SendHttpRequest("PUT", requestUrl, requestBody, jc.getHttpHeaders(), 200)
	return err
}

func (jc *JiraClient) DeleteWorklog(issueKey string, jiraWorklogId int) error {
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog/%d?notifyUsers=false", Config.JiraUrl, issueKey, jiraWorklogId)
	_, err := SendHttpRequest("DELETE", requestUrl, nil, jc.getHttpHeaders(), 204)
	return err
}
