package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_, err := GetConfig("incorrect environment variable")
	assert.Equal(t, ErrEnvVarNotFound, err)

	// just random environment variable
	_, err = GetConfig("PATH")
	assert.Equal(t, ErrFileNotFound, err)

	// correct environment variable
	_, err = GetConfig("SHORTENER_CONFIG_PATH")
	assert.Nil(t, err)
}
