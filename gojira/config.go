package gojira

import (
	"fmt"
	"os"
)

func GetEnv(key string) (string, error) {
	env, found := os.LookupEnv(key)
	if !found || env == "" {
		return "", fmt.Errorf("env %s is not set - run `gojira config` for help", key)
	}
	return env, nil
}

type Configuration struct {
	JiraUrl, JiraLogin, JiraToken, TempoUrl, TempoToken, JiraAccountId string
	UpdateExistingWorklog                                              bool
}

var Config *Configuration

const (
	envJiraInstanceURL = "GOJIRA_JIRA_INSTANCE_URL"
	envJiraLogin       = "GOJIRA_JIRA_LOGIN"
	envJiraToken       = "GOJIRA_JIRA_TOKEN" //nolint:gosec
	envJiraAccountID   = "GOJIRA_JIRA_ACCOUNT_ID"
	envTempoToken      = "GOJIRA_TEMPO_TOKEN" //nolint:gosec
)

func PrepareConfig() error {
	jiraUrl, err := GetEnv(envJiraInstanceURL)
	if err != nil {
		return err
	}
	jiraLogin, err := GetEnv(envJiraLogin)
	if err != nil {
		return err
	}
	jiraToken, err := GetEnv(envJiraToken)
	if err != nil {
		return err
	}
	jiraAccountId, err := GetEnv(envJiraAccountID)
	if err != nil {
		return err
	}
	tempoToken, err := GetEnv(envTempoToken)
	if err != nil {
		return err
	}

	Config = &Configuration{
		JiraUrl:               jiraUrl,
		JiraLogin:             jiraLogin,
		JiraToken:             jiraToken,
		JiraAccountId:         jiraAccountId,
		TempoUrl:              "https://api.tempo.io/4",
		TempoToken:            tempoToken,
		UpdateExistingWorklog: true,
	}
	return nil
}
