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
	podInfo, err := k8s.GetPodInfo()
	if err == nil {
		log.Info("K8s Pod environment detected")
		Config.Address = podInfo.GetPodIP()
		Config.Port = podInfo.GetDefaultPort()
		director, err := podInfo.GetDirector()
		if err == nil {
			Config.Director = director
		} else {
			log.Errorf("could not find VaaS director in Pod info: %s", err)
		}
	} else {
		log.Errorf("K8s Pod not detected: %s", err)
	}

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
			Name:        action.FlagDirector,
			Usage:       "VaaS director to register this backend with",
			Destination: &Config.Director,
			Value:       Config.Director,
		},
		cli.StringFlag{
			Name:        action.FlagAddress,
			Usage:       "IP address of this backend",
			Destination: &Config.Address,
			Value:       Config.Address,
		},
		cli.IntFlag{
			Name:        action.FlagPort,
			Usage:       "port of this backend",
			Destination: &Config.Port,
			Value:       Config.Port,
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
