package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	pb "github.com/siyopao/ipcheck/api/proto/v1"
	"github.com/siyopao/ipcheck/blocklist"
	"google.golang.org/grpc"
)

type appConfig struct {
	serverAddr string
}

func main() {
	appConfig, blConfig := loadEnvVars()

	if err := blocklist.InitBlocklists(blConfig); err != nil {
		checkError(err, "error initializing the blocklists")
	}

	s := grpc.NewServer()
	pb.RegisterIpCheckServer(s, &server{})

	lis, err := net.Listen("tcp", appConfig.serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadEnvVars() (appConfig, blocklist.BlConfig) {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading '.env'; continuing anyway")
	}

	serverAddr, ok := os.LookupEnv("SERVER_ADDR")
	checkOk(ok, "error reading 'SERVER_ADDR' environment variable")

	ipSetsDir, ok := os.LookupEnv("IP_SETS_DIR")
	checkOk(ok, "error reading 'IP_SETS_DIR' environment variable")

	ipSetsString, ok := os.LookupEnv("IP_SETS")
	checkOk(ok, "error reading 'IP_SETS' environment variable")
	ipSets := strings.Split(ipSetsString, ",")

	return appConfig{serverAddr}, blocklist.BlConfig{
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
