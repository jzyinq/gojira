package gojira

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// UserConfig holds user-customizable settings loaded from ~/.config/gojira/config.yaml.
type UserConfig struct {
	// SubtractedIssues lists issue keys (e.g. "PPURLOP-4") whose logged time should be
	// subtracted from the total time spent, instead of being added to it.
	SubtractedIssues []string `yaml:"subtractedIssues"`
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

// ResolveSubtractedIssueIDs resolves configured issue keys to the numeric Jira issue IDs
// that worklogs are keyed by, so they can be matched during time spent calculations.
func ResolveSubtractedIssueIDs(issueKeys []string) map[int]bool {
	issueIDs := map[int]bool{}
	for _, issueKey := range issueKeys {
		issue, err := app.jiraClient.GetIssue(issueKey)
		if err != nil {
			logrus.Warnf("failed to resolve subtracted issue %q: %v", issueKey, err)
			continue
		}
		issueIDs[issue.GetIdAsInt()] = true
	}
	return issueIDs
}
