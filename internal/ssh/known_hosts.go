package ssh

import (
	"bufio"
	"os"
	"strings"

	"github.com/recode-sh/cli/internal/system"
)

const KnownHostsFilePerm os.FileMode = 0644

type KnownHosts struct {
	knownHostsFilePath string
}

func NewKnownHosts(knownHostsFilePath string) KnownHosts {
	return KnownHosts{
		knownHostsFilePath: knownHostsFilePath,
	}
}

func NewKnownHostsWithDefaultKnownHostsFilePath() KnownHosts {
	return NewKnownHosts(
		system.DefaultSSHKnownHostsFilePath(),
	)
}

func (k KnownHosts) AddOrReplace(hostname, algorithm, fingerprint string) error {
	f, err := k.openFile()

	if err != nil {
		return err
	}

	defer f.Close()

	knownHostToAdd := hostname + " " + algorithm + " " + fingerprint
	knownHostToAddReplaced := false

	scanner := bufio.NewScanner(f)
	newKnownHostsContent := ""

	for scanner.Scan() {
		knownHostLine := scanner.Text()

		if strings.HasPrefix(knownHostLine, hostname+" "+algorithm) {
			newKnownHostsContent += knownHostToAdd + system.NewLineChar
			knownHostToAddReplaced = true
			continue
		}

		newKnownHostsContent += knownHostLine + system.NewLineChar
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !knownHostToAddReplaced {
		newKnownHostsContent += knownHostToAdd + system.NewLineChar
	}

	return os.WriteFile(
		k.knownHostsFilePath,
		[]byte(newKnownHostsContent),
		KnownHostsFilePerm,
	)
}

func (k KnownHosts) RemoveIfExists(hostname string) error {
	if len(hostname) == 0 {
		// Nothing to do.
		// We don't want to remove all hostnames
		// (the function "hasPrefix" will always return "true" if prefix is empty).
		// See below.
		return nil
	}

	f, err := k.openFile()

	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	newKnownHostsContent := ""

	for scanner.Scan() {
		knownHostLine := scanner.Text()

		if strings.HasPrefix(knownHostLine, hostname) {
			continue
		}

		newKnownHostsContent += knownHostLine + system.NewLineChar
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	/* We only want one new line
	   at the end of the file */

	for {
		trimedNewKnownHostsContent := strings.TrimSuffix(
			newKnownHostsContent,
			system.NewLineChar,
		)

		if trimedNewKnownHostsContent == newKnownHostsContent {
			break
		}

		newKnownHostsContent = trimedNewKnownHostsContent
	}

	newKnownHostsContent += system.NewLineChar

	return os.WriteFile(
		k.knownHostsFilePath,
		[]byte(newKnownHostsContent),
		KnownHostsFilePerm,
	)
}

func (k KnownHosts) openFile() (*os.File, error) {
	// create the "known_hosts" file if necessary
	return os.OpenFile(
		k.knownHostsFilePath,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		KnownHostsFilePerm,
	)
}
