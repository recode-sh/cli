package interfaces

import (
	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/recode/github"
)

type UserConfigManager interface {
	GetString(key config.UserConfigKey) string
	GetBool(key config.UserConfigKey) bool
	Set(key config.UserConfigKey, value interface{})
	WriteConfig() error
	PopulateFromGitHubUser(githubUser *github.AuthenticatedUser)
}
