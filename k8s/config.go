package k8s

import (
	"context"
	"fmt"

	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	log "github.com/sirupsen/logrus"
)

const (
	keyVaaSUser = "user"
	keyVaaSKey  = "key"
	keyVaaSURL  = "url"
)

// Config describes a k8s ConfigMap
type Config struct {
	*corev1.ConfigMap
}

// GetVaaSURL returns VaaS URL
func (cm Config) GetVaaSURL() string {
	url, err := cm.getValue(keyVaaSURL)
	if err != nil {
		log.Errorf("VaaS URL not found: %s", err)
	}
	return url
}

// GetVaaSUser returns VaaS API Username
func (cm Config) GetVaaSUser() string {
	username, err := cm.getValue(keyVaaSUser)
	if err != nil {
		log.Errorf("VaaS Username not found: %s", err)
	}
	return username
}

// GetVaaSKey returns VaaS API secret key
func (cm Config) GetVaaSKey() string {
	key, err := cm.getValue(keyVaaSKey)
	if err != nil {
		log.Errorf("VaaS secret key not found: %s", err)
	}
	return key
}

func (cm Config) getValue(key string) (string, error) {
	data := cm.GetData()
	if data[key] != "" {
		return data[key], nil
	}
	return "", fmt.Errorf("field %q not found", key)
}

// GetVaaSConfig fetches k8s PodInfo for the current Pod
func GetVaaSConfig() (*Config, error) {
	ctx := context.Background()
	podClient, err := clientProvider()
	if err != nil {
		return nil, err
	}

	pod, err := podClient.GetConfigMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get VaaS configMap: $%s", err)
	}

	return &Config{pod}, err
}
