package gojira

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("returns value for existing env var", func(t *testing.T) {
		key := "TEST_GOJIRA_ENV_VAR"
		expectedValue := "test-value"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		value, err := GetEnv(key)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("returns error for missing env var", func(t *testing.T) {
		key := "NON_EXISTENT_ENV_VAR"
		os.Unsetenv(key) // ensure it doesn't exist

		value, err := GetEnv(key)
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "env NON_EXISTENT_ENV_VAR is not set")
	})

	t.Run("returns error for empty env var", func(t *testing.T) {
		key := "EMPTY_ENV_VAR"
		os.Setenv(key, "")
		defer os.Unsetenv(key)

		value, err := GetEnv(key)
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "env EMPTY_ENV_VAR is not set")
	})
}

func TestPrepareConfig(t *testing.T) {
	t.Run("successfully prepares config with all env vars", func(t *testing.T) {
		// Setup environment variables
		os.Setenv("GOJIRA_JIRA_INSTANCE_URL", "https://test.atlassian.net")
		os.Setenv("GOJIRA_JIRA_LOGIN", "test@example.com")
		os.Setenv("GOJIRA_JIRA_TOKEN", "test-jira-token")
		os.Setenv("GOJIRA_JIRA_ACCOUNT_ID", "test-account-id")
		os.Setenv("GOJIRA_TEMPO_TOKEN", "test-tempo-token")

		defer func() {
			os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
			os.Unsetenv("GOJIRA_JIRA_LOGIN")
			os.Unsetenv("GOJIRA_JIRA_TOKEN")
			os.Unsetenv("GOJIRA_JIRA_ACCOUNT_ID")
			os.Unsetenv("GOJIRA_TEMPO_TOKEN")
			Config = nil
		}()

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
		os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
		defer func() { Config = nil }()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_INSTANCE_URL")
	})

	t.Run("returns error when GOJIRA_JIRA_LOGIN is missing", func(t *testing.T) {
		os.Setenv("GOJIRA_JIRA_INSTANCE_URL", "https://test.atlassian.net")
		os.Unsetenv("GOJIRA_JIRA_LOGIN")

		defer func() {
			os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
			Config = nil
		}()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_LOGIN")
	})

	t.Run("returns error when GOJIRA_JIRA_TOKEN is missing", func(t *testing.T) {
		os.Setenv("GOJIRA_JIRA_INSTANCE_URL", "https://test.atlassian.net")
		os.Setenv("GOJIRA_JIRA_LOGIN", "test@example.com")
		os.Unsetenv("GOJIRA_JIRA_TOKEN")

		defer func() {
			os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
			os.Unsetenv("GOJIRA_JIRA_LOGIN")
			Config = nil
		}()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_TOKEN")
	})

	t.Run("returns error when GOJIRA_JIRA_ACCOUNT_ID is missing", func(t *testing.T) {
		os.Setenv("GOJIRA_JIRA_INSTANCE_URL", "https://test.atlassian.net")
		os.Setenv("GOJIRA_JIRA_LOGIN", "test@example.com")
		os.Setenv("GOJIRA_JIRA_TOKEN", "test-token")
		os.Unsetenv("GOJIRA_JIRA_ACCOUNT_ID")

		defer func() {
			os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
			os.Unsetenv("GOJIRA_JIRA_LOGIN")
			os.Unsetenv("GOJIRA_JIRA_TOKEN")
			Config = nil
		}()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_JIRA_ACCOUNT_ID")
	})

	t.Run("returns error when GOJIRA_TEMPO_TOKEN is missing", func(t *testing.T) {
		os.Setenv("GOJIRA_JIRA_INSTANCE_URL", "https://test.atlassian.net")
		os.Setenv("GOJIRA_JIRA_LOGIN", "test@example.com")
		os.Setenv("GOJIRA_JIRA_TOKEN", "test-token")
		os.Setenv("GOJIRA_JIRA_ACCOUNT_ID", "test-account-id")
		os.Unsetenv("GOJIRA_TEMPO_TOKEN")

		defer func() {
			os.Unsetenv("GOJIRA_JIRA_INSTANCE_URL")
			os.Unsetenv("GOJIRA_JIRA_LOGIN")
			os.Unsetenv("GOJIRA_JIRA_TOKEN")
			os.Unsetenv("GOJIRA_JIRA_ACCOUNT_ID")
			Config = nil
		}()

		err := PrepareConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GOJIRA_TEMPO_TOKEN")
	})
}
