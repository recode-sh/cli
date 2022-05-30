package config

var (
	GitHubOAuthClientID    = "a8c368bfe297f0b1808a"
	GitHubOAuthCLIToAPIURL = "http://127.0.0.1:8080/github/oauth/callback"

	GitHubOAuthAPIToCLIURLPath = "/github/oauth/callback"

	GitHubOAuthScopes = []string{
		"read:user",
		"user:email",
		"repo",
		"admin:public_key",
		"admin:gpg_key",
	}
)
