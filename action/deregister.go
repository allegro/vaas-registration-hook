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
	// DeregisterName is the CLI name of this action
	DeregisterName = "deregister"
	// FlagBackendID represents a known backend id that is to be deregistered
	FlagBackendID = "backend-id, id"
)

// DeregisterCLI removes a backend from VaaS using CLI data
func DeregisterCLI(c *cli.Context) error {
	config := getCommonParameters(c.Parent().Parent())

	if config.Director == "" {
		return errors.New("no VaaS director specified")
	}

	err := config.GetSecretFromFile(config.VaaSKeyFile)
	if err != nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)
	backendID := c.Int(FlagBackendID)
	if backendID == 0 {
		bid, err := apiClient.FindBackendID(config.Director, config.Address, config.Port)
		if err != nil {
			return fmt.Errorf("could not determine backend ID: %s", err)
		}
		backendID = bid
	}

	if backendID != 0 {
		if err := apiClient.DeleteBackend(backendID); err != nil {
			return fmt.Errorf("could not deregister: %s", err)
		}

		log.WithField(FlagBackendID, backendID).
			Info("Successfully scheduled backend for deletion via VaaS")
		return nil
	}
	return errors.New("backend ID not provided")
}

// DeregisterK8s configures a VaaS client from K8s data and removes a backend
func DeregisterK8s(podInfo *k8s.PodInfo, config CommonConfig) (err error) {
	config.Address = podInfo.GetPodIP()
	config.Port = podInfo.GetDefaultPort()
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
	if err == nil {
		return fmt.Errorf("error reading VaaS secret key: %s", err)
	}

	apiClient := vaas.NewClient(config.VaaSURL, config.VaaSUser, config.VaaSKey)

	backendID, err := apiClient.FindBackendID(config.Director, config.Address, config.Port)
	if err != nil {
		return fmt.Errorf("could not determine backend ID: %s", err)
	}
	log.Infof("Deregistering backend %d from director %s", backendID, config.Director)
	if err := apiClient.DeleteBackend(backendID); err != nil {
		return fmt.Errorf("could not deregister: %s", err)
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
