package k8s

import (
	"context"
	"fmt"
	"strconv"

	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

// Annotation keys
const (
	keyDC       = "podDC"
	keyEnv      = "podEnvironment"
	keyDirector = "podDirector"
	keyWeight   = "podWeight"
	keyVaaSUser = "vaasUser"
	keyVaaSKey  = "vaasKey"
	keyVaaSURL  = "vaasUrl"
)

// PodInfo describes a k8s Pod
type PodInfo struct {
	*corev1.Pod
}

// FindAnnotation looks up an Annotation by key
func (pi PodInfo) FindAnnotation(lookupKey string) string {
	annotations := pi.GetMetadata().GetAnnotations()

	for key, value := range annotations {
		if key == lookupKey {
			return value
		}
	}

	// annotation can be nonexistent or empty
	return ""
}

// GetPorts returns a Pod's ports
func (pi PodInfo) GetPorts() []*int32 {
	containers := pi.GetSpec().GetContainers()
	// TODO(tz) Allow to specify which containers and ports will be registered
	if len(containers) > 0 {
		return pi.getPorts(containers[0])
	}

	return []*int32{}
}

// GetDefaultPort returns the first available port
func (pi PodInfo) GetDefaultPort() int {
	ports := pi.GetPorts()

	if len(ports) > 0 {
		return int(*ports[0])
	}

	return 0
}

func (pi PodInfo) getPorts(container *corev1.Container) []*int32 {
	var ports []*int32
	for _, port := range container.Ports {
		ports = append(ports, port.HostPort)
	}

	return ports
}

// GetWeight returns a Pods Weight
func (pi PodInfo) GetWeight() (int, error) {
	weight := pi.FindAnnotation(keyWeight)
	if weight == "" {
		return 0, fmt.Errorf("weight annotation is empty, annotation key: %s", keyWeight)
	}
	return strconv.Atoi(weight)
}

// GetDataCenter returns a Pod's datacenter
func (pi PodInfo) GetDataCenter() (string, error) {
	dc := pi.FindAnnotation(keyDC)
	if dc == "" {
		return "", fmt.Errorf("dc annotation is empty, annotation key: %s", keyDC)
	}
	return dc, nil
}

// GetEnvironment returns a Pod's dev/test/prod environment
func (pi PodInfo) GetEnvironment() (string, error) {
	environment := pi.FindAnnotation(keyEnv)
	if environment == "" {
		return "", fmt.Errorf("environment annotation is empty, annotation key: %s", keyEnv)
	}
	return environment, nil
}

// GetVaaSURL returns VaaS URL
func (pi PodInfo) GetVaaSURL() string {
	url := pi.FindAnnotation(keyVaaSURL)
	return url
}

// GetVaaSUser returns VaaS API Username
func (pi PodInfo) GetVaaSUser() string {
	username := pi.FindAnnotation(keyVaaSUser)
	return username
}

// GetVaaSKey returns VaaS API Username
func (pi PodInfo) GetVaaSKey() string {
	key := pi.FindAnnotation(keyVaaSKey)
	return key
}

// GetPodIP returns a Pod IP address
func (pi PodInfo) GetPodIP() string {
	return pi.GetStatus().GetPodIP()
}

// GetDirector looks up VaaS director name in Pod annotations
func (pi PodInfo) GetDirector() (string, error) {
	director := pi.FindAnnotation(keyDirector)
	if director == "" {
		return "", fmt.Errorf("director annotation is empty, annotation key: %s", keyDirector)
	}
	return director, nil
}

// GetPodInfo fetches k8s PodInfo for the current Pod
func GetPodInfo() (*PodInfo, error) {
	ctx := context.Background()
	podClient, err := clientProvider()
	if err != nil {
		return nil, err
	}

	pod, err := podClient.GetPod(ctx)
	if err != nil {
		return nil, err
	}

	return &PodInfo{pod}, err
}
