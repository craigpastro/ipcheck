package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	ipSetsDir = "/tmp/ipsets/"
)

func checkError(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v\n", msg, err)
	}
}

func main() {
	err := godotenv.Load()
	checkError("Error loading `.env`", err)

	initDb()
	defer dbPool.Close()

	r := setupRouter()
	r.Run(os.Getenv("SERVER_ADDR"))
}
