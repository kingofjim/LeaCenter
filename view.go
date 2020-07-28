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
	region := c.Param("region")
	var version, md5 string
	if versionType == "web" || versionType == "proxy" || versionType == "dns" {

		if versionType == "web" {
			if region == "cn" {
				version = serial_cn_web
			} else if region == "global" {
				version = serial_global_web
			}
		} else if versionType == "proxy" {
			if region == "cn" {
				version = serial_cn_proxy
			} else if region == "global" {
				version = serial_global_proxy
			}
		} else if versionType == "dns" {
			if region == "global" {
				version = serial_global_dns
			}
		}
		if version != "" {
			md5 = Md5sum(region, versionType)
			if md5 != "" {
				c.Data(200, "text/plain; charset=utf-8", []byte("version:"+version+"\nmd5sum:"+md5))
			}
		} else {
			c.AbortWithStatus(500)
		}
	} else {
		c.AbortWithStatus(400)
	}
}

func
Download(c *gin.Context) {
	versionType := c.Param("type")
	region := c.Param("region")
	version := c.Param("version")
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=leadns.tar.gz")
	c.Header("Content-Type", "application/gzip")
	//c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write(GetFile(versionType,region, version))
	fmt.Println(versionType, version)
}
