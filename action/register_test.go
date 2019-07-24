package action

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIfOverrides(t *testing.T) {
	value := "value1"

	result, err := overrideValue(value, "value2", "value")

	require.NoError(t, err)
	require.Equal(t, "value2", result)
}

func TestIfDoesNotOverrideOnEmpty(t *testing.T) {
	value := "value1"

	result, err := overrideValue(value, "", "value")

	require.NoError(t, err)
	require.Equal(t, "value1", result)
}

func TestIfOverrideErrorsOnEmpty(t *testing.T) {
	value := ""

	result, err := overrideValue(value, "", "value")

	require.Error(t, err)
	require.Equal(t, "no value for value", err.Error())
	require.Equal(t, "", result)
}
