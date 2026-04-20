package gojira

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// gojiraEnvVars lists all env vars used by PrepareConfig.
var gojiraEnvVars = []string{
	"GOJIRA_JIRA_INSTANCE_URL",
	"GOJIRA_JIRA_LOGIN",
	"GOJIRA_JIRA_TOKEN",
	"GOJIRA_JIRA_ACCOUNT_ID",
	"GOJIRA_TEMPO_TOKEN",
}

// isolateGojiraEnv clears all GOJIRA env vars for the duration of the test,
// then sets only those provided in set.
func isolateGojiraEnv(t *testing.T, set map[string]string) {
	t.Helper()
	for _, key := range gojiraEnvVars {
		orig, exists := os.LookupEnv(key)
		if err := os.Unsetenv(key); err == nil && exists {
			t.Cleanup(func() { os.Setenv(key, orig) }) //nolint:errcheck,gosec
		}
	}
	for k, v := range set {
		t.Setenv(k, v)
	}
}

func TestGetEnv(t *testing.T) {
	t.Run("returns value for existing env var", func(t *testing.T) {
		key := "TEST_GOJIRA_ENV_VAR"
		t.Setenv(key, "test-value")

		value, err := GetEnv(key)
		assert.NoError(t, err)
		assert.Equal(t, "test-value", value)
	})

	t.Run("returns error for missing env var", func(t *testing.T) {
		value, err := GetEnv("NON_EXISTENT_ENV_VAR")
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "env NON_EXISTENT_ENV_VAR is not set")
	})

	t.Run("returns error for empty env var", func(t *testing.T) {
		key := "EMPTY_ENV_VAR"
		t.Setenv(key, "")

		value, err := GetEnv(key)
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "env EMPTY_ENV_VAR is not set")
	})
}

func TestPrepareConfig(t *testing.T) { //nolint:funlen
	t.Run("successfully prepares config with all env vars", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			"GOJIRA_JIRA_INSTANCE_URL": "https://test.atlassian.net",
			"GOJIRA_JIRA_LOGIN":        "test@example.com",
			"GOJIRA_JIRA_TOKEN":        "test-jira-token",
			"GOJIRA_JIRA_ACCOUNT_ID":   "test-account-id",
			"GOJIRA_TEMPO_TOKEN":       "test-tempo-token",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.NoError(t, err)
		assert.NotNil(t, Config)
		assert.Equal(t, "https://test.atlassian.net", Config.JiraUrl)
		assert.Equal(t, "test@example.com", Config.JiraLogin)
		assert.Equal(t, "test-jira-token", Config.JiraToken)
		assert.Equal(t, "test-account-id", Config.JiraAccountId)
		assert.Equal(t, "test-tempo-token", Config.TempoToken)
		assert.Equal(t, "https://api.tempo.io/4", Config.TempoUrl)
		assert.True(t, Config.UpdateExistingWorklog)
	})

	t.Run("returns error when GOJIRA_JIRA_INSTANCE_URL is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_INSTANCE_URL")
	})

	t.Run("returns error when GOJIRA_JIRA_LOGIN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			"GOJIRA_JIRA_INSTANCE_URL": "https://test.atlassian.net",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_LOGIN")
	})

	t.Run("returns error when GOJIRA_JIRA_TOKEN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			"GOJIRA_JIRA_INSTANCE_URL": "https://test.atlassian.net",
			"GOJIRA_JIRA_LOGIN":        "test@example.com",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_TOKEN")
	})

	t.Run("returns error when GOJIRA_JIRA_ACCOUNT_ID is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			"GOJIRA_JIRA_INSTANCE_URL": "https://test.atlassian.net",
			"GOJIRA_JIRA_LOGIN":        "test@example.com",
			"GOJIRA_JIRA_TOKEN":        "test-token",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_ACCOUNT_ID")
	})

	t.Run("returns error when GOJIRA_TEMPO_TOKEN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			"GOJIRA_JIRA_INSTANCE_URL": "https://test.atlassian.net",
			"GOJIRA_JIRA_LOGIN":        "test@example.com",
			"GOJIRA_JIRA_TOKEN":        "test-token",
			"GOJIRA_JIRA_ACCOUNT_ID":   "test-account-id",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_TEMPO_TOKEN")
	})
}
