package action

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/allegro/vaas-registration-hook/k8s"
	"github.com/allegro/vaas-registration-hook/vaas"
)

const (
	// RegisterName is the CLI name of this action
	RegisterName = "register"

	// FlagWeight is the weight that a backend in VaaS is initially assigned
	FlagWeight = "weight"
	// FlagDC represents datacenter short name as defined in VaaS
	FlagDC = "dc"

	canaryTag = "canary"
)

// GetRegisterFlags returns a list of flags available for this action
func GetRegisterFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  FlagWeight,
			Usage: "initial weight of this backend",
			Value: 1,
		},
		cli.StringFlag{
			Name:  FlagDC,
			Usage: "datacenter short name as defined in VaaS",
		},
	}
}

// RegisterCLI configures a VaaS client from CLI data and runs register()
func RegisterCLI(c *cli.Context) error {
	config := getCommonParameters(c.Parent().Parent())

	if config.Director == "" {
		return errors.New("no VaaS director specified")
	}
	err := config.GetSecretFromFile(config.VaaSKeyFile)
	if err == nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	weight := c.Int(FlagWeight)
	dcName := c.String(FlagDC)

	return register(apiClient, config, weight, dcName)
}

// RegisterK8s configures a VaaS client from K8s data and runs register()
func RegisterK8s(podInfo *k8s.PodInfo, config CommonConfig) error {
	config.Address = podInfo.GetPodIP()
	config.Port = podInfo.GetDefaultPort()
	director, err := podInfo.GetDirector()
	if err == nil {
		config.Director = director
	} else {
		return fmt.Errorf("could not find VaaS director in Pod info: %s", err)
	}
	config.VaaSURL = podInfo.GetVaaSURL()
	config.VaaSUser = podInfo.GetVaaSUser()
	err = config.GetSecretFromFile(config.VaaSKeyFile)
	if err == nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	weight, err := podInfo.GetWeight()
	if err != nil {
		log.Errorf("unusable weight found %q", weight)
		weight = 1
	}
	dcName, err := podInfo.GetDataCenter()
	if err != nil {
		log.Errorf("unusable DC name found %q", weight)
		weight = 1
	}

	return register(apiClient, config, weight, dcName)
}

// register adds a backend to VaaS
func register(client vaas.Client, cfg CommonConfig, weight int, dcName string) (err error) {
	var tags []string
	if cfg.Canary {
		tags = []string{canaryTag}
	}

	dc, err := client.GetDC(dcName)
	if err != nil {
		return fmt.Errorf("failed getting DC info: %s", err)
	}

	director, err := client.FindDirector(cfg.Director)

	backend := vaas.Backend{
		ID:                 nil,
		Address:            cfg.Address,
		DirectorURL:        director.ResourceURI,
		DC:                 *dc,
		Port:               cfg.Port,
		InheritTimeProfile: false,
		Weight:             &weight,
		Tags:               tags,
		ResourceURI:        "",
	}
	backendID, err := client.AddBackend(&backend)

	if err == nil {
		log.Infof("Received VaaS backend id: %s", backendID)
	}

	return
}
