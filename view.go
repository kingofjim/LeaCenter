package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetSerial(c *gin.Context) {
	versionType := c.Param("type")
	if versionType == "web" {
		c.Data(200, "text/plain; charset=utf-8", []byte(serial_web))
	} else if versionType == "proxy" {
		c.Data(200, "text/plain; charset=utf-8", []byte(serial_proxy))
	} else if versionType == "dns" {
		c.Data(200, "text/plain; charset=utf-8", []byte(serial_dns))
	} else {
		c.AbortWithStatus(400)
	}
}

func Download(c *gin.Context) {
	versionType := c.Param("type")
	version := c.Param("version")
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=leadns.tar.gz")
	c.Header("Content-Type", "application/gzip")
	//c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write(GetFile(versionType, version))
	fmt.Println(versionType, version)
}