package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	if ginMode, ok := os.LookupEnv("GIN_MODE"); ok {
		gin.SetMode(ginMode)
	}

	r := gin.Default()
	r.GET("/v1/addresses/:ipaddress", inBlocklist)

	// Just for testing.
	r.PUT("/addresses", func(c *gin.Context) {
		updateBlocklist() // TODO: Handle the error.
		c.Status(http.StatusOK)
	})

	return r
}

// What additional information should I return in the response?
type Response struct {
	IPAddress string `json:"ipAddress"`
	IsBlocked bool   `json:"isBlocked"`
}

func inBlocklist(c *gin.Context) {
	ipAddress := net.ParseIP(c.Param("ipaddress"))

	if ipAddress == nil {
		c.Status(http.StatusBadRequest)
	} else {
		log.Printf("checking blocklist for '%v'\n", ipAddress)

		ok := checkBlocklist(ipAddress)
		if !ok {
			log.Printf("'%v' is NOT in the blocklist\n", ipAddress)
			c.Status(http.StatusNoContent)
		} else {
			log.Printf("'%v' is in the blocklist\n", ipAddress)
			c.JSON(http.StatusOK, Response{ipAddress.String(), true})
		}
	}
}
