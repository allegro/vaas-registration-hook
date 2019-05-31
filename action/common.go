package action

import (
	"fmt"
	"io/ioutil"
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
	FlagUser = "user"
	// EnvVaaSUser represents the user name for Auth
	EnvVaaSUser = "VAAS_USER"
	// FlagSecretKey client key for Auth
	FlagSecretKey = "key"
	// EnvVaaSKey client key for Auth
	EnvVaaSKey = "VAAS_KEY"
	// FlagSecretKeyFile client key for Auth
	FlagSecretKeyFile = "key-file"
	// EnvVaaSKeyFile client key for Auth
	EnvVaaSKeyFile = "VAAS_KEY_FILE"
	// FlagDirector represents the director name
	FlagDirector = "director"
	// FlagAddr address of this backend
	FlagAddress = "addr"
	// FlagPort represents the port of this backend
	FlagPort = "port"
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
	VaaSKeyFile  string
	Port         int
	AsyncTimeout time.Duration
}

func getCommonParameters(c *cli.Context) CommonConfig {
	return CommonConfig{
		Debug:       c.Bool(FlagDebug),
		VaaSURL:     c.String(FlagVaaSURL),
		VaaSUser:    c.String(FlagUser),
		VaaSKeyFile: c.String(FlagSecretKeyFile),
		VaaSKey:     c.String(FlagSecretKey),
		Director:    c.String(FlagDirector),
		Address:     c.String(FlagAddress),
		Port:        c.Int(FlagPort),
		Canary:      c.Bool(FlagCanaryTag),
	}
}

// GetSecretFromFile reads a value from provided file
func (config CommonConfig) GetSecretFromFile(secretFile string) error {
	if len(secretFile) != 0 {
		secret, err := ioutil.ReadFile(secretFile)
		if err != nil {
			return fmt.Errorf("unable to read secret from file: %s, %s", secretFile, err)
		}
		config.VaaSKey = string(secret)
		return nil
	}
	return nil
}
