package utils

import (
	"os"
	"testing"

	"nabd/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDatabase(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_nabd_*.db")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	err = utils.InitDatabase(tempFile.Name())
	assert.NoError(t, err)
	_, err = os.Stat(tempFile.Name())
	assert.NoError(t, err)
}

func TestInitDatabase_InvalidPath(t *testing.T) {
	err := utils.InitDatabase("/invalid/path/that/does/not/exist/test.db")
	assert.Error(t, err)
}