package router

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siyopao/ipcheck/blocklist"
)

type Response struct {
	IP          net.IP `json:"ip"`
	InBlocklist bool   `json:"inBlocklist"`
}

func InitRouter(ginMode string, config blocklist.BlConfig) *gin.Engine {
	gin.SetMode(ginMode)
	r := gin.Default()

	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	r.PUT("/v1/addresses", func(c *gin.Context) {
		if err := blocklist.CloneRepoAndPopulateTrie(config); err != nil {
			log.Fatalf("error updating blocklist: %v", err)
		}
		c.Status(http.StatusOK)
	})

	return r
}

func inBlocklist(c *gin.Context) {
	ip := net.ParseIP(c.Param("ipaddress"))

	if ip == nil {
		c.Status(http.StatusBadRequest)
	} else {
		isBlocked, err := blocklist.InBlocklist(ip)
		if err != nil {
			log.Printf("error checking if '%v' is in the blocklists: %v", ip, err)
			c.Status(http.StatusInternalServerError)
		} else if isBlocked {
			log.Printf("'%v' is in the blocklists\n", ip)
			c.JSON(http.StatusOK, Response{ip, true})
		} else {
			log.Printf("'%v' is NOT in the blocklists\n", ip)
			c.JSON(http.StatusOK, Response{ip, false})
		}
	}
}
