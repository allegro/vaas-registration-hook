package main

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	// AppName is the name of this software
	AppName   = "vaas-registration-hook"
	flagDebug = "debug"
)

var (
	// Version holds the version of this software
	Version string
	// Config contains configuration obtained from various sources
	Config CommonConfig

	debug bool
	app   *cli.App
)

// CommonConfig represents common flag values
type CommonConfig struct {
	Debug    bool
	Director string
	Address  string
	VaaSURL  string
	VaaSUser string
	VaaSKey  string
	Port     int
}

func init() {
	// ensure we will always have logs in logfmt format
	formatter := &log.TextFormatter{
		DisableColors:    false,
		QuoteEmptyFields: true,
	}
	log.SetFormatter(formatter)
	log.Printf("Initializing %s %s", AppName, Version)

	Config = CommonConfig{}

	app = cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.HideVersion = true
	app.Usage = "Binary hook for (de)registering in VaaS."
	app.Flags = getCommonFlags()
	app.Commands = getCommands()
	sort.Sort(cli.CommandsByName(app.Commands))
}

func main() {
	app.Action = func(c *cli.Context) error {
		log.Println("No action specified, exiting. For usage see --help.")
		return nil
	}

	app.Before = func(c *cli.Context) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getCommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        flagDebug,
			Usage:       "turn on debugging output",
			Destination: &Config.Debug,
			EnvVar:      "VAAS_HOOK_DEBUG",
		},
	}
}

func getCommands() []cli.Command {
	return []cli.Command{}
}
