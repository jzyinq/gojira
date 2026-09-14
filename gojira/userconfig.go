package gojira

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// UserConfig holds user-customizable settings loaded from ~/.config/gojira/config.yaml.
type UserConfig struct {
	// ExcludedIssues lists issue keys (e.g. "PPURLOP-4") whose logged time should be
	// excluded from the total time spent entirely - neither added nor subtracted.
	ExcludedIssues []string `yaml:"excludedIssues"`
}

// UserConfigPath returns the path to the user config file, following the OS config dir
// convention (e.g. respects XDG_CONFIG_HOME on Linux, falling back to ~/.config).
func UserConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gojira", "config.yaml"), nil
}

// LoadUserConfig reads the user config file. A missing file is not an error - an empty
// UserConfig is returned so gojira behaves as it did before this file existed.
func LoadUserConfig() (*UserConfig, error) {
	path, err := UserConfigPath()
	if err != nil {
		return &UserConfig{}, err
	}

	data, err := os.ReadFile(path) //nolint:gosec // path is derived from os.UserConfigDir(), not user input
	if os.IsNotExist(err) {
		return &UserConfig{}, nil
	}
	if err != nil {
		return &UserConfig{}, err
	}

	var userConfig UserConfig
	if err := yaml.Unmarshal(data, &userConfig); err != nil {
		return &UserConfig{}, err
	}
	return &userConfig, nil
}

// ResolveExcludedIssueIDs resolves configured issue keys to the numeric Jira issue IDs
// that worklogs are keyed by, so they can be matched during time spent calculations.
func ResolveExcludedIssueIDs(issueKeys []string) map[int]bool {
	issueIDs := map[int]bool{}
	for _, issueKey := range issueKeys {
		issue, err := app.jiraClient.GetIssue(issueKey)
		if err != nil {
			logrus.Warnf("failed to resolve excluded issue %q: %v", issueKey, err)
			continue
		}
		id := issue.GetIdAsInt()
		logrus.Debugf("resolved excluded issue %q to numeric id %d", issueKey, id)
		issueIDs[id] = true
	}
	return issueIDs
}
