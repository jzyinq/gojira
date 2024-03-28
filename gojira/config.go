package gojira

import (
	"fmt"
	"os"
)

func GetEnv(key string) (env string) {
	env, found := os.LookupEnv(key)
	if !found || env == "" {
		fmt.Printf("env %s is not set - run `gojira config` for help\n", key)
		os.Exit(1)
	}
	return env
}

type Configuration struct {
	JiraUrl, JiraLogin, JiraToken, TempoUrl, TempoToken, JiraAccountId string
	UpdateExistingWorklog                                              bool
}

var Config *Configuration

func PrepareConfig() {
	Config = &Configuration{
		JiraUrl:               GetEnv("GOJIRA_JIRA_INSTANCE_URL"),
		JiraLogin:             GetEnv("GOJIRA_JIRA_LOGIN"),
		JiraToken:             GetEnv("GOJIRA_JIRA_TOKEN"),
		JiraAccountId:         GetEnv("GOJIRA_JIRA_ACCOUNT_ID"),
		TempoUrl:              "https://api.tempo.io/core/3",
		TempoToken:            GetEnv("GOJIRA_TEMPO_TOKEN"),
		UpdateExistingWorklog: true,
	}
}
