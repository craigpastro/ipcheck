package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVarsLoadsDotEnv(t *testing.T) {
	appConfig, dbConfig := loadEnvVars()

	assert.Equal(t, appConfig.serverAddr, "[::1]:50051")
	assert.Equal(t, dbConfig.IPSetsDir, "/tmp/ipsets")
	assert.Equal(t, dbConfig.IPSets[0], "feodo.ipset")
}

func TestLoadEnvVarsDoesNotOverrideSetEnvVars(t *testing.T) {
	os.Setenv("SERVER_ADDR", "localhost:9999")
	os.Setenv("GIN_MODE", "release")

	appConfig, _ := loadEnvVars()

	assert.Equal(t, appConfig.serverAddr, "localhost:9999")
}
