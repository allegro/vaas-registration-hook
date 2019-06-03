package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
)

func TestIfFailsIfKubernetesAPIFails(t *testing.T) {
	client := &MockClient{}
	client.client.On("GetPod", context.Background(), "", "").
		Return(nil, errors.New("error")).Once()

	clientProvider = func() (Client, error) {
		return client, nil
	}

	podInfo, err := GetPodInfo()

	require.Error(t, err)
	require.Empty(t, podInfo)
}

func TestIfReturnsPodInfoButWithEmptyFields(t *testing.T) {
	pod := testPod()

	client := &MockClient{}
	client.client.On("GetPod", context.Background(), "", "").
		Return(pod, nil).Once()

	clientProvider = func() (Client, error) {
		return client, nil
	}

	podInfo, err := GetPodInfo()

	require.NoError(t, err)
	require.Empty(t, podInfo.Metadata.Annotations)
	require.Empty(t, podInfo.Metadata.Labels)
	require.Equal(t, len(podInfo.Spec.Containers), 0)
}

func TestIfReturnsCorrectAnnotationData(t *testing.T) {
	pod := testPod()

	expectedVipName := "lbaas-vip"
	expectedDC := "dc0"
	expectedEnvironment := "dev"
	inputWeight := "1"
	expectedWeight, err := strconv.Atoi(inputWeight)
	require.NoError(t, err)

	pod.Metadata.Annotations = map[string]string{
		keyDirector: expectedVipName,
		keyDC:       expectedDC,
		keyEnv:      expectedEnvironment,
		keyWeight:   inputWeight,
	}

	client := &MockClient{}
	client.client.On("GetPod", context.Background(), "", "").
		Return(pod, nil).Once()

	clientProvider = func() (Client, error) {
		return client, nil
	}

	podInfo, err := GetPodInfo()
	require.NoError(t, err)

	actualDirector, _ := podInfo.GetDirector()
	require.Equal(t, expectedVipName, actualDirector)

	actualDC, _ := podInfo.GetDataCenter()
	require.Equal(t, expectedDC, actualDC)

	actualWeight, _ := podInfo.GetWeight()
	require.Equal(t, expectedWeight, actualWeight)

	actualEnvironment, _ := podInfo.GetEnvironment()
	require.Equal(t, expectedEnvironment, actualEnvironment)
}

func TestIfReturnsErroredAnnotationData(t *testing.T) {
	pod := testPod()

	client := &MockClient{}
	client.client.On("GetPod", context.Background(), "", "").
		Return(pod, nil).Once()

	clientProvider = func() (Client, error) {
		return client, nil
	}

	podInfo, err := GetPodInfo()
	require.NoError(t, err)

	_, err = podInfo.GetDirector()
	require.EqualError(t, err, fmt.Sprintf("director annotation is empty, annotation key: %s", keyDirector))

	_, err = podInfo.GetDataCenter()
	require.EqualError(t, err, fmt.Sprintf("dc annotation is empty, annotation key: %s", keyDC))

	_, err = podInfo.GetWeight()
	require.EqualError(t, err, fmt.Sprintf("weight annotation is empty, annotation key: %s", keyWeight))

	_, err = podInfo.GetEnvironment()
	require.EqualError(t, err, fmt.Sprintf("environment annotation is empty, annotation key: %s", keyEnv))

}

func testPod() *corev1.Pod {
	return &corev1.Pod{
		Spec:   &corev1.PodSpec{},
		Status: &corev1.PodStatus{},
		Metadata: &metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}
}

type MockClient struct {
	client    mock.Mock
	k8sClient mock.Mock
}

func (c *MockClient) GetPod(ctx context.Context) (*corev1.Pod, error) {
	podNamespace := os.Getenv(podNamespaceEnvVar)
	podName := os.Getenv(podNameEnvVar)

	args := c.client.Called(ctx, podNamespace, podName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corev1.Pod), args.Error(1)
}
