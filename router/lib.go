package router

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/siyopao/ipcheck/storage"
)

type BlockedIP struct {
	IPAddress  net.IP      `json:"ipAddress"`
	Blocklists []Blocklist `json:"blockLists"`
}

type Blocklist struct {
	Filename       string `json:"filename"`
	SourceFileDate string `json:"sourceFileDate"`
}

func InitRouter(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	r := gin.Default()

	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	// For testing purposes.
	r.PUT("/addresses", func(c *gin.Context) {
		if err := storage.CloneAndUpdateBlocklists(); err != nil {
			log.Printf("error updating blocklist: %v", err)
		}
		c.Status(http.StatusOK)
	})

	return r
}

func inBlocklist(c *gin.Context) {
	ipAddress := net.ParseIP(c.Param("ipaddress"))

	allMatches, err := strconv.ParseBool(c.DefaultQuery("all-matches", "false"))
	if err != nil {
		allMatches = false
	}

	if ipAddress == nil {
		c.Status(http.StatusBadRequest)
	} else {
		log.Printf("checking blocklist for '%v'\n", ipAddress)

		blockedIP, err := storage.IsIPAddressInBlocklist(ipAddress, allMatches)
		if err != nil {
			log.Printf("error checking if '%v' is in the blocklists: %v", ipAddress, err)
			c.Status(http.StatusInternalServerError)
		} else if blockedIP != nil {
			log.Printf("'%v' is in the blocklists\n", ipAddress)
			c.JSON(http.StatusOK, convertToProtocolObject(blockedIP))
		} else {
			log.Printf("'%v' is NOT in the blocklists\n", ipAddress)
			c.Status(http.StatusNoContent)
		}
	}
}

func convertToProtocolObject(blockedIP *storage.BlockedIP) BlockedIP {
	var blocklists []Blocklist
	for _, blocklist := range blockedIP.Blocklists {
		blocklists = append(blocklists, Blocklist{blocklist.Filename, blocklist.SourceFileDate})
	}
	return BlockedIP{blockedIP.Address, blocklists}
}
