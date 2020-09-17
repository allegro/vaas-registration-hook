package action

import (
	"errors"
	"fmt"
	"os"

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
	// EnvDC Environment var containing datacenter short name as defined in VaaS
	EnvDC = "CLOUD_DC"
	// InstanceFormat represents a backend instance tag
	InstanceFormat = "instance:%s_%d"

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
			Name:   FlagDC,
			Usage:  "datacenter short name as defined in VaaS",
			EnvVar: EnvDC,
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
	if err != nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	weight := c.Int(FlagWeight)
	dcName := c.String(FlagDC)

	return register(apiClient, config, weight, dcName, []string{})
}

// RegisterK8s configures a VaaS client from K8s data and runs register()
func RegisterK8s(podInfo *k8s.PodInfo, config CommonConfig) (err error) {
	config.Address = podInfo.GetPodIP()
	config.Port = podInfo.GetDefaultPort()
	config.Canary = config.Canary || podInfo.FindAnnotation("canary")

	config.Director, err = overrideValue(config.Director, podInfo.GetDirector(), "Director")
	if err != nil {
		return
	}
	config.VaaSURL, err = overrideValue(config.VaaSURL, podInfo.GetVaaSURL(), "VaaS URL")
	if err != nil {
		return
	}
	config.VaaSUser, err = overrideValue(config.VaaSUser, podInfo.GetVaaSUser(), "VaaS User")
	if err != nil {
		return
	}

	err = config.GetSecretFromFile(config.VaaSKeyFile)
	if err != nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	weight, err := podInfo.GetWeight()
	if err != nil {
		log.Errorf("unusable weight %q found: %s", weight, err)
		weight = 1
	}

	dcName := os.Getenv(EnvDC)
	podDC, err := podInfo.GetDataCenter()
	if err != nil {
		log.Errorf("unusable DC name found %q: %s", dcName, err)
	}
	dcName, err = overrideValue(dcName, podDC, "DC")
	if err != nil {
		return
	}

	tags := []string{
		createInstanceTag(podInfo),
	}
	return register(apiClient, config, weight, dcName, tags)
}

func createInstanceTag(info *k8s.PodInfo) string {
	return fmt.Sprintf(InstanceFormat, info.GetName(), info.GetDefaultPort())
}

func overrideValue(oldValue, override, name string) (string, error) {
	if override != "" {
		log.Debugf("Overriding %s (%q) with %q from podInfo", name, oldValue, override)
		return override, nil
	}
	if oldValue == "" {
		return oldValue, fmt.Errorf("no value for %s", name)
	}
	return oldValue, nil
}

// register adds a backend to VaaS
func register(client vaas.Client, cfg CommonConfig, weight int, dcName string, tags []string) (err error) {
	if cfg.Canary {
		tags = append(tags, canaryTag)
	}

	dc, err := client.GetDC(dcName)
	if err != nil {
		return fmt.Errorf("failed getting DC info: %s", err)
	}

	director, err := client.FindDirector(cfg.Director)
	if err != nil {
		return fmt.Errorf("failed finding Director: %s", err)
	}

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
	log.Infof("Adding address %q port %d to director %q (%d)", cfg.Address, cfg.Port, director.Name, director.ID)
	backendID, err := client.AddBackend(&backend, director)

	if err == nil {
		log.Infof("Received VaaS backend id: %s", backendID)
	}

	return
}
