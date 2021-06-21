package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVarsLoadsDotEnv(t *testing.T) {
	appConfig, dbConfig := loadEnvVars()

	assert.Equal(t, appConfig.ginMode, "debug")
	assert.Equal(t, appConfig.serverAddr, "127.0.0.1:8080")
	assert.Equal(t, dbConfig.DatabaseURL, "postgres://postgres:password@127.0.0.1:5432/postgres")
	assert.Equal(t, dbConfig.AllMatches, false)
	assert.Equal(t, dbConfig.IpSetsDir, "/tmp/ipsets")
	assert.Equal(t, dbConfig.IpSets[0], "feodo.ipset")
}

func TestLoadEnvVarsDoesNotOverrideSetEnvVars(t *testing.T) {
	os.Setenv("SERVER_ADDR", "localhost:9999")
	os.Setenv("GIN_MODE", "release")

	appConfig, _ := loadEnvVars()

	assert.Equal(t, appConfig.serverAddr, "localhost:9999")
	assert.Equal(t, appConfig.ginMode, "release")
}
