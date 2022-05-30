package system

import "os"

type EnvVars struct{}

func NewEnvVars() EnvVars {
	return EnvVars{}
}

func (EnvVars) Get(key string) string {
	return os.Getenv(key)
}
