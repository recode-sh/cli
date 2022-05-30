package ssh_test

import (
	"os"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
	"github.com/recode-sh/cli/internal/system"
)

func TestKnownHostsAddOrReplaceWithNonEmptyKnownHostsFile(t *testing.T) {
	knownHostsPath := "./testdata/non_empty_known_hosts"
	knownHostsAtStart, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified known hosts file
		err = os.WriteFile(
			knownHostsPath,
			knownHostsAtStart,
			os.FileMode(ssh.KnownHostsFilePerm),
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedKnownHosts := string(knownHostsAtStart) +
		"hostname algorithm fingerprint" +
		system.NewLineChar

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err = knownHosts.AddOrReplace("hostname", "algorithm", "fingerprint")

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

func TestKnownHostsAddOrReplaceWithEmptyKnownHostsFile(t *testing.T) {
	knownHostsPath := "./testdata/empty_known_hosts"
	expectedKnownHosts :=
		"hostname algorithm fingerprint" +
			system.NewLineChar

	defer func() { // Reset modified known hosts file
		err := os.WriteFile(
			knownHostsPath,
			[]byte(""),
			os.FileMode(ssh.KnownHostsFilePerm),
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err := knownHosts.AddOrReplace("hostname", "algorithm", "fingerprint")

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

func TestKnownHostsAddOrReplaceWithNonExistingKnownHostsFile(t *testing.T) {
	knownHostsPath := "./testdata/non_existing_known_hosts"
	expectedKnownHosts :=
		"hostname algorithm fingerprint" +
			system.NewLineChar

	defer func() { // Remove created known hosts file
		err := os.Remove(knownHostsPath)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	knownHosts := ssh.NewKnownHosts(knownHostsPath)
	err := knownHosts.AddOrReplace("hostname", "algorithm", "fingerprint")

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

	createdKnownHostsFileInfo, err := os.Stat(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if createdKnownHostsFileInfo.Mode().Perm() != ssh.KnownHostsFilePerm {
		t.Fatalf(
			"expected created known hosts file to have permission '%o', got '%o'",
			ssh.KnownHostsFilePerm,
			createdKnownHostsFileInfo.Mode().Perm(),
		)
	}
}

func TestKnownHostsAddOrReplaceWitExistingHost(t *testing.T) {
	knownHostsPath := "./testdata/non_empty_known_hosts"
	knownHostsAtStart, err := os.ReadFile(knownHostsPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Reset modified known hosts file
		err = os.WriteFile(
			knownHostsPath,
			knownHostsAtStart,
			os.FileMode(ssh.KnownHostsFilePerm),
		)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	expectedKnownHosts := `github.com,140.82.118.3 ssh-rsa fingerprint
34.229.126.51 ssh-rsa fingerprint

github.com ecdsa-sha2-nistp256 fingerprint
github.com ssh-ed25519 fingerprint

34.229.126.51 ecdsa-sha2-nistp256 fingerprint_replaced
34.229.126.51 ssh-ed25519 fingerprint
`

	knownHosts := ssh.NewKnownHosts(knownHostsPath)

	err = knownHosts.AddOrReplace(
		"34.229.126.51",
		"ecdsa-sha2-nistp256",
		"fingerprint_replaced",
	)

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
