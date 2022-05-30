package ssh

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/recode-sh/cli/internal/system"
)

const PrivateKeyFilePerm os.FileMode = 0600

type Keys struct {
	sshDir string
}

func NewKeys(SSHDir string) Keys {
	return Keys{
		sshDir: SSHDir,
	}
}

func NewKeysWithDefaultDir() Keys {
	return NewKeys(
		system.DefaultSSHDir(),
	)
}

func (k Keys) CreateOrReplacePEM(
	PEMName string,
	PEMContent string,
) (pathWritten string, err error) {

	pathWritten = filepath.Join(k.sshDir, PEMName+".pem")

	err = os.WriteFile(
		pathWritten,
		[]byte(PEMContent),
		PrivateKeyFilePerm,
	)

	return
}

func (k Keys) RemovePEMIfExists(PEMName string) error {
	err := os.Remove(
		k.GetPEMFilePath(PEMName),
	)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func (k Keys) GetPEMFilePath(PEMName string) string {
	return filepath.Join(
		k.sshDir,
		PEMName+".pem",
	)
}
