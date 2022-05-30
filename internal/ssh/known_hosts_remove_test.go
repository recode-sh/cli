package ssh_test

import (
	"os"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
)

func TestKnownHostsRemoveWithExistingHostname(t *testing.T) {
	knownHostsPath := "./testdata/non_empty_known_hosts"
	knownHostsAtStart, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified known hosts file
		err = os.WriteFile(
			knownHostsPath,
			knownHostsAtStart,
			ssh.KnownHostsFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedKnownHosts := `github.com,140.82.118.3 ssh-rsa fingerprint

github.com ecdsa-sha2-nistp256 fingerprint
github.com ssh-ed25519 fingerprint
`

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err = knownHosts.RemoveIfExists("34.229.126.51")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	knownHostsAtEnd, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if string(knownHostsAtEnd) != expectedKnownHosts {
		t.Fatalf(
			"expected known hosts to equal '%s', got '%s'",
			expectedKnownHosts,
			string(knownHostsAtEnd),
		)
	}
}

func TestKnownHostsRemoveWithNonExistingHostname(t *testing.T) {
	knownHostsPath := "./testdata/non_empty_known_hosts"
	knownHostsAtStart, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified known hosts file
		err = os.WriteFile(
			knownHostsPath,
			knownHostsAtStart,
			ssh.KnownHostsFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedKnownHosts := knownHostsAtStart

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err = knownHosts.RemoveIfExists("104.78.1.4")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	knownHostsAtEnd, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if string(expectedKnownHosts) != string(knownHostsAtEnd) {
		t.Fatalf(
			"expected known hosts to equal '%s', got '%s'",
			string(expectedKnownHosts),
			string(knownHostsAtEnd),
		)
	}
}

func TestKnownHostsRemoveWithEmptyHostname(t *testing.T) {
	knownHostsPath := "./testdata/non_empty_known_hosts"
	knownHostsAtStart, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified known hosts file
		err = os.WriteFile(
			knownHostsPath,
			knownHostsAtStart,
			ssh.KnownHostsFilePerm,
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedKnownHosts := knownHostsAtStart

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err = knownHosts.RemoveIfExists("")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	knownHostsAtEnd, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if string(expectedKnownHosts) != string(knownHostsAtEnd) {
		t.Fatalf(
			"expected known hosts to equal '%s', got '%s'",
			string(expectedKnownHosts),
			string(knownHostsAtEnd),
		)
	}
}
