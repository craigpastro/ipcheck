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

	if err = godotenv.Load(); err != nil {
		log.Println("error loading '.env'; continuing anyway")
	}

	serverAddr, ok = os.LookupEnv("SERVER_ADDR")
	checkOk(ok, "error reading 'SERVER_ADDR' environment variable")

	ginMode, ok = os.LookupEnv("GIN_MODE")
	checkOk(ok, "error reading 'GIN_MODE' environment variable")

	databaseURL, ok = os.LookupEnv("DATABASE_URL")
	checkOk(ok, "error reading 'DATABASE_URL' environment variable")

	allMatchesString, ok := os.LookupEnv("ALL_MATCHES")
	checkOk(ok, "error reading 'ALL_MATCHES' environment variable")
	allMatches, err = strconv.ParseBool(allMatchesString)
	checkError(err, "error parsing 'ALL_MATCHES' environment variable to bool")

	ipSetsDir, ok = os.LookupEnv("IP_SETS_DIR")
	checkOk(ok, "error reading 'IP_SETS_DIR' environment variable")

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk(ok, "error reading 'IP_SETS' environment variable")
	ipSets = strings.Split(ipSetsString, ",")
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v\n", msg, err)
	}
}

func checkOk(ok bool, msg string) {
	if !ok {
		log.Fatalf("%v\n", msg)
	}
}
