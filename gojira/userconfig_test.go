package gojira

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// isolateUserConfigDir points os.UserConfigDir() (via $XDG_CONFIG_HOME on Linux) at a
// fresh temp directory so tests never touch the real user's config file.
func isolateUserConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	return dir
}

func TestLoadUserConfig(t *testing.T) {
	t.Run("returns empty config when file is missing", func(t *testing.T) {
		isolateUserConfigDir(t)

		userConfig, err := LoadUserConfig()
		assert.NoError(t, err)
		assert.Empty(t, userConfig.ExcludedIssues)
	})

	t.Run("loads excluded issues from yaml file", func(t *testing.T) {
		dir := isolateUserConfigDir(t)
		configPath := filepath.Join(dir, "gojira", "config.yaml")
		assert.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0750))
		assert.NoError(t, os.WriteFile(configPath, []byte("excludedIssues:\n  - PPURLOP-4\n"), 0600))

		userConfig, err := LoadUserConfig()
		assert.NoError(t, err)
		assert.Equal(t, []string{"PPURLOP-4"}, userConfig.ExcludedIssues)
	})

	t.Run("returns error for invalid yaml", func(t *testing.T) {
		dir := isolateUserConfigDir(t)
		configPath := filepath.Join(dir, "gojira", "config.yaml")
		assert.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0750))
		assert.NoError(t, os.WriteFile(configPath, []byte("not: valid: yaml: ["), 0600))

		_, err := LoadUserConfig()
		assert.Error(t, err)
	})
}
