package ssh_test

import (
	"os"
	"strings"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
	"github.com/recode-sh/cli/internal/system"
)

func TestConfigUpdateHostWithExistingHost(t *testing.T) {
	configPath := "./testdata/non_empty_ssh_config"
	configAtStart, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified config file
		err = os.WriteFile(
			configPath,
			configAtStart,
			ssh.ConfigFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedConfig := `Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_rsa
  IdentitiesOnly yes
  ServerAliveInterval 240

Host 34.128.204.12
  HostName updated_hostname
  IdentityFile updated_identityFile
  User updated_user
  ForwardAgent yes`

	config := ssh.NewConfig(configPath)

	updatedHostName := "updated_hostname"
	updatedUser := "updated_user"
	updatedIdentityFile := "updated_identityFile"

	err = config.UpdateHost(
		"34.128.204.12",
		&updatedHostName,
		&updatedIdentityFile,
		&updatedUser,
	)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEnd, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEndString := strings.TrimSuffix(
		string(configAtEnd),
		system.NewLineChar,
	)

	if configAtEndString != expectedConfig {
		t.Fatalf(
			"expected config to equal '%s', got '%s'",
			expectedConfig,
			configAtEndString,
		)
	}
}

func TestConfigUpdateHostWithExistingHostAndPartialConfig(t *testing.T) {
	configPath := "./testdata/non_empty_ssh_config"
	configAtStart, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified config file
		err = os.WriteFile(
			configPath,
			configAtStart,
			ssh.ConfigFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedConfig := `Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_rsa
  IdentitiesOnly yes
  ServerAliveInterval 240

Host 34.128.204.12
  HostName updated_hostname
  IdentityFile identityFile
  User user
  ForwardAgent yes`

	config := ssh.NewConfig(configPath)

	updatedHostName := "updated_hostname"

	err = config.UpdateHost(
		"34.128.204.12",
		&updatedHostName,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEnd, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEndString := strings.TrimSuffix(
		string(configAtEnd),
		system.NewLineChar,
	)

	if configAtEndString != expectedConfig {
		t.Fatalf(
			"expected config to equal '%s', got '%s'",
			expectedConfig,
			configAtEndString,
		)
	}
}

func TestConfigUpdateHostWithNonExistingHost(t *testing.T) {
	configPath := "./testdata/non_empty_ssh_config"
	configAtStart, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified config file
		err = os.WriteFile(
			configPath,
			configAtStart,
			ssh.ConfigFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedConfig := string(configAtStart)

	config := ssh.NewConfig(configPath)

	updatedHostName := "updated_hostname"
	updatedUser := "updated_user"
	updatedIdentityFile := "updated_identityFile"

	err = config.UpdateHost(
		"34.228.204.12",
		&updatedHostName,
		&updatedIdentityFile,
		&updatedUser,
	)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEnd, err := os.ReadFile(configPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	configAtEndString := strings.TrimSuffix(
		string(configAtEnd),
		system.NewLineChar,
	)

	if configAtEndString != expectedConfig {
		t.Fatalf(
			"expected config to equal '%s', got '%s'",
			expectedConfig,
			configAtEndString,
		)
	}
}
