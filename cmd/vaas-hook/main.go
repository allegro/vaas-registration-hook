package main

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/allegro/vaas-registration-hook/action"
	"github.com/allegro/vaas-registration-hook/k8s"
)

const (
	// AppName is the name of this software
	AppName = "vaas-registration-hook"
)

var (
	// Version holds the version of this software
	Version string
	// Config contains configuration obtained from various sources
	Config action.CommonConfig

	app *cli.App
)

func init() {
	// ensure we will always have logs in logfmt format
	formatter := &log.TextFormatter{
		DisableColors:    false,
		QuoteEmptyFields: true,
	}
	log.SetFormatter(formatter)
	log.Printf("Initializing %s %s", AppName, Version)

	Config = action.CommonConfig{}

	app = cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.HideVersion = false
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
		if Config.Debug {
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
			Name:        action.FlagDebug,
			Usage:       "turn on debugging output",
			Destination: &Config.Debug,
			EnvVar:      action.EnvDebug,
		},
		cli.StringFlag{
			Name:        action.FlagVaaSURL,
			Usage:       "address of the VaaS endpoint",
			Destination: &Config.VaaSURL,
			EnvVar:      action.EnvVaaSURL,
		},
		cli.StringFlag{
			Name:        action.FlagUser,
			Usage:       "user for Auth",
			Destination: &Config.VaaSUser,
			EnvVar:      action.EnvVaaSUser,
		},
		cli.StringFlag{
			Name:        action.FlagSecretKey,
			Usage:       "client key for Auth",
			Destination: &Config.VaaSKey,
			EnvVar:      action.EnvVaaSKey,
		},
		cli.StringFlag{
			Name:        action.FlagDirector,
			Usage:       "VaaS director to register this backend with",
			Destination: &Config.Director,
		},
		cli.StringFlag{
			Name:        action.FlagAddress,
			Usage:       "IP address of this backend",
			Destination: &Config.Address,
		},
		cli.IntFlag{
			Name:        action.FlagPort,
			Usage:       "port of this backend",
			Destination: &Config.Port,
		},
		cli.BoolFlag{
			Name:        action.FlagCanaryTag,
			Destination: &Config.Canary,
			Usage:       "this backend is a canary",
		},
	}
}

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  action.RegisterName,
			Usage: "register a backend with VaaS",
			Subcommands: []cli.Command{
				{
					Name:  "cli",
					Usage: "register using data from command line/env",
					Action: func(c *cli.Context) error {
						log.Print("Registering services using data from command line/env")
						return action.RegisterCLI(c)
					},
					Flags: action.GetRegisterFlags(),
				},
				{
					Name:  "k8s",
					Usage: "register using data from Kubernetes API",
					Action: func(c *cli.Context) error {
						log.Print("Registering services using data from Kubernetes API")

						podInfo, err := k8s.GetPodInfo()
						if err != nil {
							log.Errorf("K8s Pod not detected: %s", err)
							return nil
						}
						log.Info("K8s Pod environment detected")

						vaasConfig, err := k8s.GetVaaSConfig()
						Config.VaaSURL = vaasConfig.GetVaaSURL()
						Config.VaaSUser = vaasConfig.GetVaaSUser()
						Config.VaaSKey = vaasConfig.GetVaaSKey()

						return action.RegisterK8s(podInfo, Config)
					},
				},
			},
		},
		{
			Name:  action.DeregisterName,
			Usage: "deregister a backend from VaaS",
			Subcommands: []cli.Command{
				{
					Name:  "cli",
					Usage: "Deregister using data from command line/env",
					Action: func(c *cli.Context) error {
						log.Print("Deregistering services using data from command line/env")
						return action.DeregisterCLI(c)
					},
					Flags: action.GetDeregisterFlags(),
				},
				{
					Name:  "k8s",
					Usage: "Deregister using data from Kubernetes API",
					Action: func(c *cli.Context) error {
						log.Print("Deregistering services using data from Kubernetes API")

						podInfo, err := k8s.GetPodInfo()
						if err != nil {
							log.Errorf("K8s Pod not detected: %s", err)
							return nil
						}
						log.Info("K8s Pod environment detected")

						vaasConfig, err := k8s.GetVaaSConfig()
						Config.VaaSURL = vaasConfig.GetVaaSURL()
						Config.VaaSUser = vaasConfig.GetVaaSUser()
						Config.VaaSKey = vaasConfig.GetVaaSKey()

						return action.DeregisterK8s(podInfo, Config)
					},
				},
			},
		},
	}
}
