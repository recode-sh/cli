//go:build prod

package config

func init() {
	GitHubOAuthClientID = "7e1b6c93f4ba81819162"
	GitHubOAuthCLIToAPIURL = "https://recode-sh-api.herokuapp.com/github/oauth/callback"
}
