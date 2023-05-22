package gojira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func getJiraAuthorizationHeader() string {
	authorizationToken := fmt.Sprintf("%s:%s", Config.JiraLogin, Config.JiraToken)
	authorizationHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authorizationToken)))
	return authorizationHeader
}

func (issue Issue) NewWorkLog(timeSpent string) (WorkLog, error) {
	payload := map[string]string{
		"timeSpent":      timeSpent,
		"adjustEstimate": "leave",
	}
	payloadJson, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(payloadJson)
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", Config.JiraUrl, issue.Key)
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}

	response, err := SendHttpRequest("POST", requestUrl, requestBody, headers, 201)
	if err != nil {
		return WorkLog{}, err
	}
	// print string version of response
	//fmt.Println(string(response))

	// FIXME that's a fiction that we're having actual WorkLog object here - response is different
	var worklog WorkLog
	err = json.Unmarshal(response, &worklog)
	if err != nil {
		return WorkLog{}, err
	}
	fmt.Printf("Successfully logged %s of time to ticket %s\n", timeSpent, issue.Key)
	worklog.StartDate = app.time.Format(dateLayout)
	return worklog, nil
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
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
	} `json:"fields"`
}

func GetLatestIssues() (JQLResponse, error) {
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
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}
	response, err := SendHttpRequest("POST", requestUrl, requestBody, headers, 200)
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

func GetIssue(issueKey string) (Issue, error) {
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=summary,status", Config.JiraUrl, issueKey)
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}
	response, err := SendHttpRequest("GET", requestUrl, nil, headers, 200)
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
