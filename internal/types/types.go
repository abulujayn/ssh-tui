package types

// SSHHost represents a parsed SSH host with its configuration
type SSHHost struct {
	Name     string
	HostName string
	User     string
	Port     string
	Source   string
	Aliases  []string
}

const (
	SourceConfig     = "config"
	SourceKnownHosts = "known_hosts"
	SourceCustom     = "custom"
	DefaultSSHPort   = "22"
)
