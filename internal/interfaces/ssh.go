package interfaces

// SSH known hosts

type SSHKnownHostsManager interface {
	RemoveIfExists(hostname string) error
	AddOrReplace(hostname, algorithm, fingerprint string) error
}

// SSH Config

type SSHConfigManager interface {
	AddOrReplaceHost(hostKey, hostName, identityFile, user string, port int64) error
	UpdateHost(hostKey string, hostName, identityFile, user *string) error
	RemoveHostIfExists(hostKey string) error
}

// SSH keys

type SSHKeysManager interface {
	CreateOrReplacePEM(PEMName, PEMContent string) (string, error)
	RemovePEMIfExists(PEMPath string) error
	GetPEMFilePath(PEMName string) string
}
