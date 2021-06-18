package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environment variables
var (
	serverAddr  string
	ginMode     string
	databaseURL string
	allMatches  bool
	ipSetsDir   string
	ipSets      []string
)

func main() {
	loadEnvVars()

	initDb()
	defer dbPool.Close()

	r := setupRouter()
	r.Run(os.Getenv("SERVER_ADDR"))
}

func loadEnvVars() {
	var ok bool
	var err error

	err = godotenv.Load()
	checkError("Error loading '.env'", err)

	serverAddr, ok = os.LookupEnv("SERVER_ADDR")
	checkOk("error reading 'SERVER_ADDR' environment variable", ok)

	ginMode, ok = os.LookupEnv("GIN_MODE")
	checkOk("error reading 'GIN_MODE' environment variable", ok)

	databaseURL, ok = os.LookupEnv("DATABASE_URL")
	checkOk("error reading 'DATABASE_URL' environment variable", ok)

	allMatchesString, ok := os.LookupEnv("ALL_MATCHES")
	checkOk("error reading 'ALL_MATCHES' environment variable", ok)
	allMatches, err = strconv.ParseBool(allMatchesString)
	checkError("error parsing 'ALL_MATCHES' environment variable to bool", err)

	ipSetsDir, ok = os.LookupEnv("IP_SETS_DIR")
	checkOk("error reading 'IP_SETS_DIR' environment variable", ok)

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk("error reading 'IP_SETS' environment variable", ok)
	ipSets = strings.Split(ipSetsString, ",")
}

func checkError(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v\n", msg, err)
	}
}

func checkOk(msg string, ok bool) {
	if !ok {
		log.Fatalf("%v\n", msg)
	}
}
