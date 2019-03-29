package action

import (
	"time"

	"github.com/urfave/cli"
)

// These Flag* consts exist to make any changes to flags consistent across the project
const (
	// FlagDebug turn on debugging output
	FlagDebug = "debug"
	// EnvDebug turn on debugging output
	EnvDebug = "DEBUG"
	// FlagVaaSURL address of the VaaS host to query
	FlagVaaSURL = "vaas-url"
	// EnvVaaSURL address of the VaaS host to query
	EnvVaaSURL = "VAAS_URL"
	// FlagUser represents the user name for Auth
	FlagUser = "user, u"
	// EnvVaaSUser represents the user name for Auth
	EnvVaaSUser = "VAAS_USER"
	// FlagSecretKey client key for Auth
	FlagSecretKey = "key, k"
	// EnvVaaSKey client key for Auth
	EnvVaaSKey = "VAAS_KEY"
	// FlagDirector represents the director name
	FlagDirector = "director"
	// FlagAddr address of this backend
	FlagAddress = "addr"
	// FlagPort represents the port of this backend
	FlagPort = "port, p"
	// FlagCanaryTag
	FlagCanaryTag = "canary"

	// IDFileLoc file containing VaaS backend ID
	IDFileLoc = "/tmp/vaas.id"
)

// CommonConfig represents common flag values
type CommonConfig struct {
	Debug        bool
	DryRun       bool
	Canary       bool
	Director     string
	Address      string
	VaaSURL      string
	VaaSUser     string
	VaaSKey      string
	Port         int
	AsyncTimeout time.Duration
}

func getCommonParameters(c *cli.Context) CommonConfig {
	return CommonConfig{
		Debug:    c.Bool(FlagDebug),
		VaaSURL:  c.String(FlagVaaSURL),
		VaaSUser: c.String(FlagUser),
		VaaSKey:  c.String(FlagSecretKey),
		Director: c.String(FlagDirector),
		Address:  c.String(FlagAddress),
		Port:     c.Int(FlagPort),
		Canary:   c.Bool(FlagCanaryTag),
	}
}
