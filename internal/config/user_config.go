package config

import (
	"github.com/recode-sh/recode/github"
	"github.com/spf13/viper"
)

type UserConfigKey string

const (
	UserConfigKeyUserIsLoggedIn    UserConfigKey = "user_is_logged_in"
	UserConfigKeyGitHubAccessToken UserConfigKey = "github_access_token"
	UserConfigKeyGitHubUsername    UserConfigKey = "github_username"
	UserConfigKeyGitHubEmail       UserConfigKey = "github_email"
	UserConfigKeyGitHubFullName    UserConfigKey = "github_full_name"
)

type UserConfig struct{}

func NewUserConfig() UserConfig {
	return UserConfig{}
}

func (UserConfig) GetString(key UserConfigKey) string {
	return viper.GetString(string(key))
}

func (UserConfig) GetBool(key UserConfigKey) bool {
	return viper.GetBool(string(key))
}

func (UserConfig) Set(key UserConfigKey, value interface{}) {
	viper.Set(string(key), value)
}

func (u UserConfig) PopulateFromGitHubUser(githubUser *github.AuthenticatedUser) {
	u.Set(
		UserConfigKeyGitHubEmail,
		githubUser.PrimaryEmail,
	)

	u.Set(
		UserConfigKeyGitHubFullName,
		githubUser.FullName,
	)

	u.Set(
		UserConfigKeyGitHubUsername,
		githubUser.Username,
	)
}

func (UserConfig) WriteConfig() error {
	return viper.WriteConfig()
}
