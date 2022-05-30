package system

import (
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil || !os.IsNotExist(err) {
		return true
	}

	return false
}

func UserHomeDir() string {
	// Ignore errors since we only care about Windows and *nix.
	homedir, _ := os.UserHomeDir()
	return homedir
}

// UserConfigDir returns the path where
// the user config files should be stored
// following XDG Base Directory Specification.
// Ref: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func UserConfigDir() string {
	baseConfigDir := os.Getenv("XDG_CONFIG_HOME")

	if len(baseConfigDir) == 0 {
		baseConfigDir = filepath.Join(UserHomeDir(), ".config")
	}

	return filepath.Join(baseConfigDir, "recode")
}

func UserConfigFilePath() string {
	return filepath.Join(UserConfigDir(), "recode.yml")
}

func DefaultSSHDir() string {
	return filepath.Join(UserHomeDir(), ".ssh")
}

func DefaultSSHDirExists() bool {
	return PathExists(DefaultSSHDir())
}

func DefaultSSHConfigFilePath() string {
	return filepath.Join(DefaultSSHDir(), "config")
}

func DefaultSSHConfigFileExists() bool {
	return PathExists(DefaultSSHConfigFilePath())
}

func DefaultSSHKnownHostsFilePath() string {
	return filepath.Join(DefaultSSHDir(), "known_hosts")
}

func DefaultSSHKnownHostsFileExists() bool {
	return PathExists(DefaultSSHKnownHostsFilePath())
}
