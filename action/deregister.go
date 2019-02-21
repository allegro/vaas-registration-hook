package action

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/allegro/vaas-registration-hook/vaas"
)

const (
	// DeregisterName is the CLI name of this action
	DeregisterName = "deregister"
	// FlagBackendID represents a known backend id that is to be deregistered
	FlagBackendID = "backend-id, id"
)

// Deregister removes a backend from VaaS
func Deregister(c *cli.Context) error {
	config := getCommonParameters(c.Parent())
	log.Debug(config)

	apiClient := vaas.NewClient(config.VaaSUrl, config.Username, config.Key)
	backendID := c.Int(FlagBackendID)

	if backendID != 0 {
		if err := apiClient.DeleteBackend(backendID); err != nil {
			return err
		}

		log.WithField(FlagBackendID, backendID).
			Info("Successfully scheduled backend for deletion via VaaS")
	} else {
		return errors.New("backend ID not provided")
	}

	return nil
}

// GetDeregisterFlags returns a list of flags available for this action
func GetDeregisterFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  FlagBackendID,
			Usage: "known backend id that is to be deregistered",
		},
	}
}
