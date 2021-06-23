package router

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siyopao/ipcheck/blocklist"
)

type Response struct {
	InBlocklist bool `json:"inBlocklist"`
}

func InitRouter(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	r := gin.Default()

	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	// For testing purposes.
	// r.PUT("/addresses", func(c *gin.Context) {
	// 	if err := blocklist.CloneAndUpdateBlocklists(); err != nil {
	// 		log.Printf("error updating blocklist: %v", err)
	// 	}
	// 	c.Status(http.StatusOK)
	// })

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
			c.JSON(http.StatusOK, Response{true})
		} else {
			log.Printf("'%v' is NOT in the blocklists\n", ip)
			c.JSON(http.StatusOK, Response{false})
		}
	}
}
