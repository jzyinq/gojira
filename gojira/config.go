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

func PrepareConfig() error {
	jiraUrl, err := GetEnv("GOJIRA_JIRA_INSTANCE_URL")
	if err != nil {
		return err
	}
	jiraLogin, err := GetEnv("GOJIRA_JIRA_LOGIN")
	if err != nil {
		return err
	}
	jiraToken, err := GetEnv("GOJIRA_JIRA_TOKEN")
	if err != nil {
		return err
	}
	jiraAccountId, err := GetEnv("GOJIRA_JIRA_ACCOUNT_ID")
	if err != nil {
		return err
	}
	tempoToken, err := GetEnv("GOJIRA_TEMPO_TOKEN")
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
