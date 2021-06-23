package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/siyopao/ipcheck/blocklist"
	"github.com/siyopao/ipcheck/router"
)

type appConfig struct {
	serverAddr string
	ginMode    string
}

func main() {
	appConfig, blConfig := loadEnvVars()

	if err := blocklist.CloneRepoAndPopulateTrie(blConfig); err != nil {
		checkError(err, "error initializing the blocklists")
	}
	c := cron.New()
	c.AddFunc("@every 24h", func() {
		if err := blocklist.CloneRepoAndPopulateTrie(blConfig); err != nil {
			log.Printf("error updating blocklist: %v", err)
		}
	})
	c.Start()

	r := router.TestMode(router.InitRouter(appConfig.ginMode), blConfig)
	if err := r.Run(appConfig.serverAddr); err != nil {
		checkError(err, "error starting the server")
	}
}

func loadEnvVars() (appConfig, blocklist.BlConfig) {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading '.env'; continuing anyway")
	}

	serverAddr, ok := os.LookupEnv("SERVER_ADDR")
	checkOk(ok, "error reading 'SERVER_ADDR' environment variable")

	ginMode, ok := os.LookupEnv("GIN_MODE")
	checkOk(ok, "error reading 'GIN_MODE' environment variable")

	ipSetsDir, ok := os.LookupEnv("IP_SETS_DIR")
	checkOk(ok, "error reading 'IP_SETS_DIR' environment variable")

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk(ok, "error reading 'IP_SETS' environment variable")
	ipSets := strings.Split(ipSetsString, ",")

	return appConfig{serverAddr, ginMode}, blocklist.BlConfig{
		IPSetsDir: ipSetsDir,
		IPSets:    ipSets,
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
