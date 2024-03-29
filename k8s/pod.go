package k8s

import (
	"context"
	"fmt"
	"strconv"

	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

// Annotation keys
const (
	keyDC                = "podDC"
	keyEnv               = "podEnvironment"
	keyDirector          = "podDirector"
	keyWeight            = "podWeight"
	keyVaaSUser          = "vaasUser"
	keyVaaSURL           = "vaasUrl"
	SidecarContainerName = "envoy-sidecar"
)

// PodInfo describes a k8s Pod
type PodInfo struct {
	*corev1.Pod
}

// GetAnnotation looks up an Annotation by key
func (pi PodInfo) GetAnnotation(lookupKey string) string {
	annotations := pi.GetMetadata().GetAnnotations()

	for key, value := range annotations {
		if key == lookupKey {
			return value
		}
	}

	// annotation can be nonexistent or empty
	return ""
}

// FindAnnotation looks up an Annotation by key
func (pi PodInfo) FindAnnotation(lookupKey string) bool {
	annotations := pi.GetMetadata().GetAnnotations()

	for key := range annotations {
		if key == lookupKey {
			return true
		}
	}

	// annotation can be nonexistent
	return false
}

// GetPorts returns a Pod's ports
func (pi PodInfo) GetPorts() []*int32 {
	containers := pi.GetSpec().GetContainers()
	// TODO(tz) Allow to specify which containers and ports will be registered
	if len(containers) > 0 {
		for _, container := range containers {
			if container.GetName() == SidecarContainerName {
				return pi.getPorts(container)
			}
		}
		for _, container := range containers {
			if len(container.Ports) > 0 {
				return pi.getPorts(container)
			}
		}
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
		ports = append(ports, port.ContainerPort)
	}

	return ports
}

// GetWeight returns a Pods Weight
func (pi PodInfo) GetWeight() (int, error) {
	weight := pi.GetAnnotation(keyWeight)
	if weight == "" {
		return 0, fmt.Errorf("weight annotation is empty, annotation key: %s", keyWeight)
	}
	return strconv.Atoi(weight)
}

// GetDataCenter returns a Pod's datacenter
func (pi PodInfo) GetDataCenter() (string, error) {
	dc := pi.GetAnnotation(keyDC)
	if dc == "" {
		return "", fmt.Errorf("dc annotation is empty, annotation key: %s", keyDC)
	}
	return dc, nil
}

// GetEnvironment returns a Pod's dev/test/prod environment
func (pi PodInfo) GetEnvironment() (string, error) {
	environment := pi.GetAnnotation(keyEnv)
	if environment == "" {
		return "", fmt.Errorf("environment annotation is empty, annotation key: %s", keyEnv)
	}
	return environment, nil
}

// GetVaaSURL returns VaaS URL
func (pi PodInfo) GetVaaSURL() string {
	return pi.GetAnnotation(keyVaaSURL)
}

// GetVaaSUser returns VaaS API Username
func (pi PodInfo) GetVaaSUser() string {
	return pi.GetAnnotation(keyVaaSUser)
}

// GetPodIP returns a Pod IP address
func (pi PodInfo) GetPodIP() string {
	return pi.GetStatus().GetPodIP()
}

// GetDirector looks up VaaS director name in Pod annotations
func (pi PodInfo) GetDirector() string {
	return pi.GetAnnotation(keyDirector)
}

// GetUID retrieves Pod UID form it's metadata
func (pi PodInfo) GetUID() *string {
	return pi.Metadata.Uid
}

// GetName retrieves Pod name form it's metadata
func (pi PodInfo) GetName() string {
	return pi.Metadata.GetName()
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
