package main

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/allegro/vaas-registration-hook/action"
)

const (
	// AppName is the name of this software
	AppName   = "vaas-registration-hook"
	flagDebug = "debug"
)

var (
	// Version holds the version of this software
	Version string

	debug bool
	app   *cli.App
)

func init() {
	// ensure we will always have logs in logfmt format
	formatter := &log.TextFormatter{
		DisableColors:    false,
		QuoteEmptyFields: true,
	}
	log.SetFormatter(formatter)
	log.Printf("Initializing $s %s", AppName, Version)

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
			Destination: &debug,
			EnvVar:      "VAAS_HOOK_DEBUG",
		},
		cli.StringFlag{
			Name:   action.FlagVaaSURL,
			Usage:  "address of the VaaS endpoint",
			EnvVar: "VAAS_URL",
		},
		cli.StringFlag{
			Name:   action.FlagUser,
			Usage:  "user for Auth",
			EnvVar: "VAAS_USER",
		},
		cli.StringFlag{
			Name:   action.FlagSecretKey,
			Usage:  "client key for Auth",
			EnvVar: "VAAS_KEY",
		},
		cli.StringFlag{
			Name:  action.FlagDirector,
			Usage: "VaaS director to register this backend with",
		},
		cli.StringFlag{
			Name:  action.FlagAddress,
			Usage: "IP address of this backend",
		},
		cli.IntFlag{
			Name:  action.FlagPort,
			Usage: "port of this backend",
		},
		cli.BoolFlag{
			Name:  action.FlagDryRun,
			Usage: "do not perform any changes, just print requests",
		},
		cli.BoolFlag{
			Name:  action.FlagCanaryTag,
			Usage: "this backend is a canary",
		},
	}
}

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   action.RegisterName,
			Usage:  "register a backend with VaaS",
			Action: action.Register,
			Flags:  action.GetRegisterFlags(),
		},
		{
			Name:   action.DeregisterName,
			Usage:  "deregister a backend from VaaS",
			Action: action.Deregister,
			Flags:  action.GetDeregisterFlags(),
		},
	}
}
