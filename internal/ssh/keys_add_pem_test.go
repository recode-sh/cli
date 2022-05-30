package ssh_test

import (
	"os"
	"testing"

	"github.com/recode-sh/cli/internal/ssh"
)

func TestKeysCreateOrReplacePEM(t *testing.T) {
	keys := ssh.NewKeys("./testdata")

	expectedPEMName := "pem_name"
	expectedPEMContent := "pem_content"
	expectedPEMPath := "./testdata/" + expectedPEMName + ".pem"

	PEMPath, err := keys.CreateOrReplacePEM(
		expectedPEMName,
		expectedPEMContent,
	)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	defer func() { // Remove created PEM file
		err = os.Remove(PEMPath)

		if err != nil {
			t.Fatalf("expected no error, got '%+v'", err)
		}
	}()

	if "./"+PEMPath != expectedPEMPath {
		t.Fatalf(
			"expected PEM path to equal '%s', got '%s'",
			expectedPEMPath,
			PEMPath,
		)
	}

	createdPEMContent, err := os.ReadFile(expectedPEMPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if string(createdPEMContent) != expectedPEMContent {
		t.Fatalf(
			"expected PEM to equal '%s', got '%s'",
			expectedPEMContent,
			string(createdPEMContent),
		)
	}

	createdPEMFileInfo, err := os.Stat(expectedPEMPath)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if createdPEMFileInfo.Mode().Perm() != ssh.PrivateKeyFilePerm {
		t.Fatalf(
			"expected created PEM file to have permission '%o', got '%o'",
			ssh.PrivateKeyFilePerm,
			createdPEMFileInfo.Mode().Perm(),
		)
	}
}
