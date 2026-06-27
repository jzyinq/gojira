package gojira

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testJiraURL   = "https://test.atlassian.net"
	testJiraEmail = "test@example.com"
)

// gojiraEnvVars lists all env vars used by PrepareConfig.
var gojiraEnvVars = []string{
	envJiraInstanceURL,
	envJiraLogin,
	envJiraToken,
	envJiraAccountID,
	envTempoToken,
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
			envJiraInstanceURL: testJiraURL,
			envJiraLogin:       testJiraEmail,
			envJiraToken:       "test-jira-token",
			envJiraAccountID:   "test-account-id",
			envTempoToken:      "test-tempo-token",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.NoError(t, err)
		assert.NotNil(t, Config)
		assert.Equal(t, testJiraURL, Config.JiraUrl)
		assert.Equal(t, testJiraEmail, Config.JiraLogin)
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
		assert.Contains(t, err.Error(), envJiraInstanceURL)
	})

	t.Run("returns error when GOJIRA_JIRA_LOGIN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			envJiraInstanceURL: testJiraURL,
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), envJiraLogin)
	})

	t.Run("returns error when GOJIRA_JIRA_TOKEN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			envJiraInstanceURL: testJiraURL,
			envJiraLogin:       testJiraEmail,
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), envJiraToken)
	})

	t.Run("returns error when GOJIRA_JIRA_ACCOUNT_ID is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			envJiraInstanceURL: testJiraURL,
			envJiraLogin:       testJiraEmail,
			envJiraToken:       "test-token",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), envJiraAccountID)
	})

	t.Run("returns error when GOJIRA_TEMPO_TOKEN is missing", func(t *testing.T) {
		isolateGojiraEnv(t, map[string]string{
			envJiraInstanceURL: testJiraURL,
			envJiraLogin:       testJiraEmail,
			envJiraToken:       "test-token",
			envJiraAccountID:   "test-account-id",
		})
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), envTempoToken)
	})
}
