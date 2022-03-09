package gojira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

func getJiraAuthorizationHeader() string {
	authorizationToken := fmt.Sprintf("%s:%s", Config.JiraLogin, Config.JiraToken)
	authorizationHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authorizationToken)))
	return authorizationHeader
}

func (issue Issue) NewWorkLog(timeSpent string) {
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

	SendHttpRequest("POST", requestUrl, requestBody, headers, 201)
	fmt.Printf("Successfully logged %s of time to ticket %s\n", timeSpent, issue.Key)
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

func GetLatestIssues() JQLResponse {
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
	response := SendHttpRequest("POST", requestUrl, requestBody, headers, 200)
	var jqlResponse JQLResponse
	err = json.Unmarshal(response, &jqlResponse)
	if err != nil {
		panic(err)
	}
	return jqlResponse
}

func GetIssue(issueKey string) Issue {
	requestUrl := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=summary,status", Config.JiraUrl, issueKey)
	headers := map[string]string{
		"Authorization": getJiraAuthorizationHeader(),
		"Content-Type":  "application/json",
	}
	response := SendHttpRequest("GET", requestUrl, nil, headers, 200)
	var jiraIssue Issue
	err := json.Unmarshal(response, &jiraIssue)
	if err != nil {
		log.Fatalln(err)
	}
	return jiraIssue
}
