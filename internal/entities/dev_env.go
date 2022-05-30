package entities

type DevEnvAdditionalProperties struct {
	GitHubCreatedSSHKeyId *int64 `json:"github_created_ssh_key_id"`
	GitHubCreatedGPGKeyId *int64 `json:"github_created_gpg_key_id"`
}
