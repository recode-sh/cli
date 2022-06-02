package system_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/recode-sh/cli/internal/system"
)

func TestPathExistsWithExistingPath(t *testing.T) {
	existingPath := "./paths_test.go"
	pathExists := system.PathExists(existingPath)

	if !pathExists {
		t.Fatalf("expected 'true', got 'false'")
	}
}

func TestPathExistsWithNonExistingPath(t *testing.T) {
	nonExistingPath := "./path-that-doesnt-exist"
	pathExists := system.PathExists(nonExistingPath)

	if pathExists {
		t.Fatalf("expected 'false', got 'true'")
	}
}

func TestUserHomeDir(t *testing.T) {
	expectedHomeDir, err := os.UserHomeDir()

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if system.UserHomeDir() != expectedHomeDir {
		t.Fatalf(
			"expected user home directory to equal '%s', got '%s'",
			expectedHomeDir,
			system.UserHomeDir(),
		)
	}
}

func TestDefaultSSHDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	expectedSSHDir := filepath.Join(homeDir, ".ssh")

	if system.DefaultSSHDir() != expectedSSHDir {
		t.Fatalf(
			"expected default SSH directory to equal '%s', got '%s'",
			expectedSSHDir,
			system.DefaultSSHDir(),
		)
	}
}

func TestDefaultSSHConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	expectedSSHConfigFilePath := filepath.Join(homeDir, ".ssh/config")

	if system.DefaultSSHConfigFilePath() != expectedSSHConfigFilePath {
		t.Fatalf(
			"expected default SSH config file path to equal '%s', got '%s'",
			expectedSSHConfigFilePath,
			system.DefaultSSHConfigFilePath(),
		)
	}
}
