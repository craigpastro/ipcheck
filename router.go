package main

import (
	"log"
	"net"
	"net/http"

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

func setupRouter(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	r := gin.Default()

	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	// For testing purposes.
	r.PUT("/addresses", func(c *gin.Context) {
		if err := updateBlocklists(); err != nil {
			log.Printf("error updating blocklist: %v", err)
		}
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
