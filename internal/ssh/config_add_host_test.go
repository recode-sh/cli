package ssh_test

import (
	"os"
	"strings"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
	"github.com/recode-sh/cli/internal/system"
)

func TestConfigAddOrReplaceHostWithNonEmptyConfig(t *testing.T) {
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

	expectedConfig := string(configAtStart) + `
Host hostkey
# added by Recode
  HostName hostname
  IdentityFile identityFile
  User user
  Port 2200
  ForwardAgent yes`

	config := ssh.NewConfig(configPath)
	err = config.AddOrReplaceHost(
		"hostkey",
		"hostname",
		"identityFile",
		"user",
		2200,
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

func TestConfigAddOrReplaceHostWithEmptyConfig(t *testing.T) {
	configPath := "./testdata/empty_ssh_config"
	expectedConfig := `Host hostkey
# added by Recode
  HostName hostname
  IdentityFile identityFile
  User user
  Port 2200
  ForwardAgent yes`

	defer func() { // Reset modified config file
		err := os.WriteFile(
			configPath,
			[]byte(""),
			ssh.ConfigFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	config := ssh.NewConfig(configPath)
	err := config.AddOrReplaceHost(
		"hostkey",
		"hostname",
		"identityFile",
		"user",
		2200,
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

func TestConfigAddOrReplaceHostWithNonExistingConfig(t *testing.T) {
	configPath := "./testdata/non_existing_ssh_config"
	expectedConfig := `Host hostkey
# added by Recode
  HostName hostname
  IdentityFile identityFile
  User user
  Port 2200
  ForwardAgent yes`

	defer func() { // Remove created config file
		err := os.Remove(configPath)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	config := ssh.NewConfig(configPath)
	err := config.AddOrReplaceHost(
		"hostkey",
		"hostname",
		"identityFile",
		"user",
		2200,
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

func TestConfigAddOrReplaceHostWithInvalidHostKey(t *testing.T) {
	configPath := "./testdata/non_empty_ssh_config"
	invalidHostKey := ""

	config := ssh.NewConfig(configPath)
	err := config.AddOrReplaceHost(
		invalidHostKey,
		"hostname",
		"identityFile",
		"user",
		2200,
	)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}
}

func TestConfigAddOrReplaceHostWithExistingHostConfig(t *testing.T) {
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
# added by Recode
  HostName hostname_replaced
  IdentityFile identityFile_replaced
  User user_replaced
  Port 2200
  ForwardAgent yes`

	config := ssh.NewConfig(configPath)
	err = config.AddOrReplaceHost(
		"34.128.204.12",
		"hostname_replaced",
		"identityFile_replaced",
		"user_replaced",
		2200,
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
