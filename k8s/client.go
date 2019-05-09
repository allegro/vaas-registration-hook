package k8s

import (
	"context"
	"fmt"
	"os"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

const (
	podNamespaceEnvVar = "KUBERNETES_POD_NAMESPACE"
	podNameEnvVar      = "KUBERNETES_POD_NAME"
	configMapEnvVar    = "KUBERNETES_CONFIG_MAP"
)

// Client contains methods to access k8s API
type Client interface {
	// GetPod returns current pod data.
	GetPod(ctx context.Context) (*corev1.Pod, error)
	GetConfigMap(ctx context.Context) (*corev1.ConfigMap, error)
}

var clientProvider = func() (Client, error) {
	k8sClient, err := k8s.NewInClusterClient()

	return &defaultClient{k8sClient: k8sClient}, err
}

type defaultClient struct {
	k8sClient *k8s.Client
}

// GetPod returns k8s Pod information
func (c *defaultClient) GetPod(ctx context.Context) (*corev1.Pod, error) {
	podNamespace := os.Getenv(podNamespaceEnvVar)
	podName := os.Getenv(podNameEnvVar)

	pod := &corev1.Pod{}
	if err := c.k8sClient.Get(ctx, podNamespace, podName, pod); err != nil {
		return nil, fmt.Errorf("unable to get pod data from API: %s", err)
	}

	return pod, nil
}

// GetConfigMap returns a k8s ConfigMap
func (c *defaultClient) GetConfigMap(ctx context.Context) (*corev1.ConfigMap, error) {
	podNamespace := os.Getenv(podNamespaceEnvVar)
	mapName := os.Getenv(configMapEnvVar)

	config := &corev1.ConfigMap{}
	if err := c.k8sClient.Get(ctx, podNamespace, mapName, config); err != nil {
		return nil, fmt.Errorf("unable to get config data from API: %s", err)
	}

	return config, nil
}
