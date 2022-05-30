package ssh

import (
	"fmt"
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/recode-sh/cli/internal/system"
)

const ConfigFilePerm os.FileMode = 0644

type Config struct {
	configFilePath string
}

func NewConfig(configFilePath string) Config {
	return Config{
		configFilePath: configFilePath,
	}
}

func NewConfigWithDefaultConfigFilePath() Config {
	return NewConfig(system.DefaultSSHConfigFilePath())
}

func (c Config) AddOrReplaceHost(
	hostKey string,
	hostName string,
	identityFile string,
	user string,
	port int64,
) error {

	cfg, err := c.parse()

	if err != nil {
		return err
	}

	hostPattern, err := ssh_config.NewPattern(hostKey)

	if err != nil {
		return err
	}

	hostNodes := []ssh_config.Node{
		&ssh_config.Empty{
			Comment: " added by Recode",
		},

		&ssh_config.KV{
			Key:   "  HostName",
			Value: hostName,
		},

		&ssh_config.KV{
			Key:   "  IdentityFile",
			Value: identityFile,
		},

		&ssh_config.KV{
			Key:   "  User",
			Value: user,
		},

		&ssh_config.KV{
			Key:   "  Port",
			Value: fmt.Sprint(port),
		},

		&ssh_config.KV{
			Key:   "  ForwardAgent",
			Value: "yes",
		},
	}

	hostToAdd := &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{
			hostPattern,
		},
		Nodes: hostNodes,
	}

	hostToAddIndex := c.lookupHostIndex(
		cfg,
		hostKey,
	)

	if hostToAddIndex == -1 {
		cfg.Hosts = append(cfg.Hosts, hostToAdd)
	} else {
		cfg.Hosts[hostToAddIndex] = hostToAdd
	}

	return c.save(cfg)
}

func (c Config) UpdateHost(
	hostKey string,
	hostName *string,
	identityFile *string,
	user *string,
) error {

	cfg, err := c.parse()

	if err != nil {
		return err
	}

	updatedHosts := []*ssh_config.Host{}

	for _, host := range cfg.Hosts {
		// We don't use "host.Matches()"
		// here because we don't want the
		// wildcard host ("Host *") to match
		if len(host.Patterns) == 1 && host.Patterns[0].String() == hostKey {
			for _, node := range host.Nodes {
				switch t := node.(type) {
				case *ssh_config.KV:
					lowercasedKey := strings.ToLower(t.Key)

					if lowercasedKey == "hostname" && hostName != nil {
						t.Value = *hostName
					}

					if lowercasedKey == "identityfile" && identityFile != nil {
						t.Value = *identityFile
					}

					if lowercasedKey == "user" && user != nil {
						t.Value = *user
					}
				}
			}
		}

		updatedHosts = append(updatedHosts, host)
	}

	cfg.Hosts = updatedHosts

	return c.save(cfg)
}

func (c Config) RemoveHostIfExists(hostKey string) error {
	cfg, err := c.parse()

	if err != nil {
		return err
	}

	updatedHosts := []*ssh_config.Host{}

	for _, host := range cfg.Hosts {
		// We don't use "host.Matches()"
		// here because we don't want the
		// wildcard host ("Host *") to match
		if len(host.Patterns) == 1 && host.Patterns[0].String() == hostKey {
			continue
		}

		updatedHosts = append(updatedHosts, host)
	}

	cfg.Hosts = updatedHosts

	return c.save(cfg)
}

func (c Config) lookupHostIndex(
	cfg *ssh_config.Config,
	hostKey string,
) int {

	for hostIndex, host := range cfg.Hosts {
		// We don't use "host.Matches()"
		// here because we don't want the
		// wildcard host ("Host *") to match
		if len(host.Patterns) == 1 && host.Patterns[0].String() == hostKey {
			return hostIndex
		}

		continue
	}

	return -1
}

func (c Config) parse() (*ssh_config.Config, error) {
	f, err := os.OpenFile(
		c.configFilePath,
		os.O_CREATE|os.O_RDONLY,
		ConfigFilePerm,
	)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ssh_config.Decode(f)
}

func (c Config) save(cfg *ssh_config.Config) error {
	return os.WriteFile(
		c.configFilePath,
		[]byte(cfg.String()),
		ConfigFilePerm,
	)
}
