package action

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSecretFromFileCorrectly(t *testing.T) {
	config := CommonConfig{
		VaaSKeyFile: "vaas-secret",
	}
	secretContent := "secret token"

	data := []byte(secretContent)
	err := ioutil.WriteFile(config.VaaSKeyFile, data, 0644)
	require.NoError(t, err)

	err = config.GetSecretFromFile(config.VaaSKeyFile)
	require.NoError(t, err)
	require.Equal(t, config.VaaSKey, secretContent)

	err = os.Remove(config.VaaSKeyFile)
	require.NoError(t, err)
}
