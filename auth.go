package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func AuthToken(token string, originToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Query("token")
		if clientToken == token {
			isOrigin = false
		} else if clientToken == originToken && CheckIP(c.ClientIP()){
			isOrigin = true
		} else {
			c.AbortWithStatus(401)
		}
	}
}

func CheckIP(clientIP string) bool {
	if clientIP == "::1" {
		clientIP = "127.0.0.1"
	}
	fmt.Println("IP: ", clientIP)
	originIP := os.Getenv("ORIGIN_IP")
	ipList := strings.Split(originIP, ",")
	for _, ip := range ipList {
		if ip == clientIP {
			return true
		}
	}
	return false
}