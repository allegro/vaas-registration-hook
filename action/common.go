package action

import (
	"time"

	"github.com/urfave/cli"
)

// These Flag* consts exist to make any changes to flags consistent across the project
const (
	// FlagVaaSURL address of the VaaS host to query
	FlagVaaSURL = "vaas-url"
	// FlagUser represents the user name for Auth
	FlagUser = "user, u"
	// FlagSecretKey client key for Auth
	FlagSecretKey = "key, k"
	// FlagDirector represents the director name
	FlagDirector = "director"
	// FlagAddr address of this backend
	FlagAddress = "addr"
	// FlagPort represents the port of this backend
	FlagPort = "port, p"
	// FlagDryRun only perform read operations, do not make any changes
	FlagDryRun = "dry-run"
	// FlagCanaryTag
	FlagCanaryTag = "canary"
	// FlagAsyncTimeout
	FlagAsyncTimeout = "timeout, t"
)

type commonParameters struct {
	VaaSUrl      string
	Username     string
	Key          string
	Director     string
	Address      string
	Port         int
	DryRun       bool
	AsyncTimeout time.Duration
	Canary       bool
}

func getCommonParameters(c *cli.Context) commonParameters {
	return commonParameters{
		VaaSUrl:      c.String(FlagVaaSURL),
		Username:     c.String(FlagUser),
		Key:          c.String(FlagSecretKey),
		Director:     c.String(FlagDirector),
		Address:      c.String(FlagAddress),
		Port:         c.Int(FlagPort),
		DryRun:       c.Bool(FlagDryRun),
		AsyncTimeout: c.Duration(FlagAsyncTimeout),
		Canary:       c.Bool(FlagCanaryTag),
	}
}
