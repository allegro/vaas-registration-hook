package action

import (
	"errors"
	"io/ioutil"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/allegro/vaas-registration-hook/k8s"
	"github.com/allegro/vaas-registration-hook/vaas"
)

const (
	// DeregisterName is the CLI name of this action
	DeregisterName = "deregister"
	// FlagBackendID represents a known backend id that is to be deregistered
	FlagBackendID = "backend-id, id"
)

// DeregisterCLI removes a backend from VaaS using CLI data
func DeregisterCLI(c *cli.Context) error {
	config := getCommonParameters(c.Parent().Parent())
	debug(config)

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	backendID := c.Int(FlagBackendID)

	if backendID != 0 {
		if err := apiClient.DeleteBackend(backendID); err != nil {
			return err
		}

		log.WithField(FlagBackendID, backendID).
			Info("Successfully scheduled backend for deletion via VaaS")
	}
	return errors.New("backend ID not provided")
}

// DeregisterK8s configures a VaaS client from K8s data and removes a backend
func DeregisterK8s(podInfo *k8s.PodInfo, config CommonConfig) error {
	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)

	backendID, err := loadBackendID()
	if err != nil {
		return err
	}
	if err := apiClient.DeleteBackend(backendID); err != nil {
		return err
	}

	return errors.New("backend ID not provided")
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

func loadBackendID() (s int, err error) {
	data, err := ioutil.ReadFile(IDFileLoc)
	if err == nil {
		s, err = strconv.Atoi(string(data))
	}
	return
}
