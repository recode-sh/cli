package ssh_test

import (
	"errors"
	"os"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
)

func TestKeysRemoveExistingPEM(t *testing.T) {
	keys := ssh.NewKeys("./testdata")

	PEMName := "pem_to_remove"
	PEMPath := "./testdata/" + PEMName + ".pem"

	_, err := os.Create(PEMPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	err = keys.RemovePEMIfExists(PEMName)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	_, err = os.Stat(PEMPath)

	if err == nil {
		t.Fatalf("expected file not exists error, got nothing")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file not exists error, got '%+v'", err)
	}
}

func TestKeysRemoveNonExistingPEM(t *testing.T) {
	keys := ssh.NewKeys("./testdata")
	err := keys.RemovePEMIfExists("non_existing_pem")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}
}
