package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

type appConfig struct {
	serverAddr string
	ginMode    string
}

type dbConfig struct {
	databaseURL string
	allMatches  bool
	ipSetsDir   string
	ipSets      []string
}

func main() {
	appConfig, dbConfig := loadEnvVars()

	initDb(dbConfig)
	defer dbPool.Close()

	// Immediately initialize the blocklists in the db. Then run approximately
	// daily after that.
	//
	// NOTE: In some sense adding and updating the blocklists is adding state
	// to this service and, instead, this should probably be done in a Lambda.
	checkError(cloneAndUpdateBlocklists(), "error initializing the blocklists")
	c := cron.New()
	c.AddFunc("@every 25h3m", func() {
		if err := cloneAndUpdateBlocklists(); err != nil {
			log.Printf("error updating blocklist: %v", err)
		}
	})
	c.Start()

	r := setupRouter(appConfig.ginMode)
	r.Run(appConfig.serverAddr)
}

func loadEnvVars() (appConfig, dbConfig) {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading '.env'; continuing anyway")
	}

	serverAddr, ok := os.LookupEnv("SERVER_ADDR")
	checkOk(ok, "error reading 'SERVER_ADDR' environment variable")

	ginMode, ok := os.LookupEnv("GIN_MODE")
	checkOk(ok, "error reading 'GIN_MODE' environment variable")

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	checkOk(ok, "error reading 'DATABASE_URL' environment variable")

	allMatchesString, ok := os.LookupEnv("ALL_MATCHES")
	checkOk(ok, "error reading 'ALL_MATCHES' environment variable")
	allMatches, err := strconv.ParseBool(allMatchesString)
	checkError(err, "error parsing 'ALL_MATCHES' environment variable to bool")

	ipSetsDir, ok := os.LookupEnv("IP_SETS_DIR")
	checkOk(ok, "error reading 'IP_SETS_DIR' environment variable")

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk(ok, "error reading 'IP_SETS' environment variable")
	ipSets := strings.Split(ipSetsString, ",")

	return appConfig{serverAddr, ginMode}, dbConfig{databaseURL, allMatches, ipSetsDir, ipSets}
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
