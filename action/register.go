package action

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

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
		},
		cli.StringFlag{
			Name:  FlagDC,
			Usage: "datacenter short name as defined in VaaS",
		},
	}
}

// Register adds a backend to VaaS
func Register(c *cli.Context) error {
	config := getCommonParameters(c.Parent())
	log.Debug(config)

	apiClient := vaas.NewClient(config.VaaSUrl, config.Username, config.Key)
	weight := c.Int(FlagWeight)
	dcName := c.String(FlagDC)

	return register(apiClient, config, weight, dcName)
}

func register(client vaas.Client, cfg commonParameters, weight int, dcName string) error {
	var tags []string
	if cfg.Canary {
		tags = []string{canaryTag}
	}

	dc, err := client.GetDC(dcName)
	if err != nil {
		return fmt.Errorf("failed getting DC info: %s", err)
	}

	backend := vaas.Backend{
		ID:                 nil,
		Address:            cfg.Address,
		Director:           cfg.Director,
		DC:                 *dc,
		Port:               cfg.Port,
		InheritTimeProfile: false,
		Weight:             &weight,
		Tags:               tags,
		ResourceURI:        "",
	}
	_, err = client.AddBackend(&backend)
	return err
}
