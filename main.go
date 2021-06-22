package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/siyopao/ipcheck/router"
	"github.com/siyopao/ipcheck/storage"
)

type appConfig struct {
	serverAddr string
	ginMode    string
}

func main() {
	appConfig, dbConfig := loadEnvVars()

	if err := storage.InitDb(dbConfig); err != nil {
		checkError(err, "error initilizing the db")
	}
	defer storage.DbPool.Close()

	// Immediately initialize the blocklists in the db. Then run approximately
	// daily after that. Could set a specific time instead, or actually have
	// this set as a configuation variable.
	if err := storage.CloneAndUpdateBlocklists(); err != nil {
		checkError(err, "error initializing the blocklists")
	}
	c := cron.New()
	c.AddFunc("@every 24h", func() {
		if err := storage.CloneAndUpdateBlocklists(); err != nil {
			log.Printf("error updating blocklist: %v", err)
		}
	})
	c.Start()

	r := router.InitRouter(appConfig.ginMode)
	if err := r.Run(appConfig.serverAddr); err != nil {
		checkError(err, "error starting the server")
	}
}

func loadEnvVars() (appConfig, storage.DbConfig) {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading '.env'; continuing anyway")
	}

	serverAddr, ok := os.LookupEnv("SERVER_ADDR")
	checkOk(ok, "error reading 'SERVER_ADDR' environment variable")

	ginMode, ok := os.LookupEnv("GIN_MODE")
	checkOk(ok, "error reading 'GIN_MODE' environment variable")

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	checkOk(ok, "error reading 'DATABASE_URL' environment variable")

	ipSetsDir, ok := os.LookupEnv("IP_SETS_DIR")
	checkOk(ok, "error reading 'IP_SETS_DIR' environment variable")

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk(ok, "error reading 'IP_SETS' environment variable")
	ipSets := strings.Split(ipSetsString, ",")

	return appConfig{serverAddr, ginMode}, storage.DbConfig{
		DatabaseURL: databaseURL,
		IPSetsDir:   ipSetsDir,
		IPSets:      ipSets,
	}
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
