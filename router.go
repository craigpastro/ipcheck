package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Response struct {
	IPAddress  net.IP      `json:"ipAddress"`
	Blocklists []Blocklist `json:"blockLists"`
}

type Blocklist struct {
	Filename       string `json:"filename"`
	SourceFileDate string `json:"sourceFileDate"`
}

func setupRouter() *gin.Engine {
	if ginMode, ok := os.LookupEnv("GIN_MODE"); ok {
		gin.SetMode(ginMode)
	}

	r := gin.Default()
	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	// Just for testing.
	r.PUT("/addresses", func(c *gin.Context) {
		go updateBlocklist() // TODO: How to handle the error?
		c.Status(http.StatusOK)
	})

	return r
}

func inBlocklist(c *gin.Context) {
	ipAddress := net.ParseIP(c.Param("ipaddress"))

	if ipAddress == nil {
		c.Status(http.StatusBadRequest)
	} else {
		log.Printf("checking blocklist for '%v'\n", ipAddress)

		blockedIP, err := isIPAddressInBlocklist(ipAddress)
		if err != nil {
			log.Printf("error checking if '%v' is in the blocklists: %v", ipAddress, err)
			c.Status(http.StatusInternalServerError)
		} else if blockedIP != nil {
			log.Printf("'%v' is in the blocklists\n", ipAddress)
			c.JSON(http.StatusOK, Response{blockedIP.address, convertToProtocolObject(blockedIP.blocklists)})
		} else {
			log.Printf("'%v' is NOT in the blocklists\n", ipAddress)
			c.Status(http.StatusNoContent)
		}
	}
}

func convertToProtocolObject(blocklists []blocklist) []Blocklist {
	var res []Blocklist
	for _, blocklist := range blocklists {
		res = append(res, Blocklist{blocklist.filename, blocklist.sourceFileDate})
	}
	return res
}
