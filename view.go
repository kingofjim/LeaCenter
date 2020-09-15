package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
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
	if versionType == "web" || versionType == "proxy" || versionType == "dns" || versionType == "all" {
		if versionType == "all" && region == "global" {
			version = serial_global_all
		} else if versionType == "web" {
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
			md5 = Md5sum(versionType, region, version)
			if md5 != "" {
				if versionType == "all" {
					c.Data(200, "text/plain; charset=utf-8", []byte("version:"+serial_global_all+"\nmd5sum:"+md5+"\ndnsVersion:"+serial_global_dns+"\nwebVersion:"+serial_global_web+"\nproxyVersion:"+serial_global_proxy))
				} else {
					c.Data(200, "text/plain; charset=utf-8", []byte("version:"+version+"\nmd5sum:"+md5))
				}
			}
		} else {
			c.AbortWithStatus(500)
		}
	} else {
		c.AbortWithStatus(400)
	}
}

func GetAllSerial(c *gin.Context) {
	fmt.Println("haha")
}

func Download(c *gin.Context) {
	versionType := c.Param("type")
	region := c.Param("region")
	version := c.Param("version")
	c.Writer.WriteHeader(http.StatusOK)
	var filename string
	if versionType == "all" {
		filename = "compressfile.tar.gz"
	} else if versionType == "dns" {
		filename = "dns.tar.gz"
	} else {
		filename = "leadns.tar.gz"
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/gzip")
	//c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write(GetFile(versionType, region, version))
	//fmt.Println(versionType, version)
}

func Commit(c *gin.Context) {
	if isOrigin == true {
		var data commitData
		var status int
		c.BindJSON(&data)
		//fmt.Println("Request:", c.Request)
		//fmt.Println("Header:", c.Request.Header)
		//fmt.Println("Body", c.Request.Body)
		log.Infoln(data)
		if data.Version == "" || data.DnsVersion == "" || data.WebVersion == "" || data.ProxyVersion == "" {
			status = 500
			c.JSON(status, gin.H{
				"error":   "value_error",
				"message": "Broken input",
			})
			return
		}

		cmd := exec.Command("tar", "zxf", "tmp/compressfile.tar.gz", "-C", "tmp")
		err := cmd.Run()

		if err != nil {
			status = 500
			c.JSON(status, gin.H{
				"error":   err.Error(),
				"message": "untar failed.",
			})
			return
		}

		fileCheck, err := CheckFileExist(&data)

		if err != nil {
			status = 500
			c.JSON(status, gin.H{
				"error": err.Error(),
			})
		} else {
			if fileCheck {
				go StoreFile(data)
				status = 200
			} else {
				status = 406
			}
			c.JSON(status, gin.H{
				"version": data.Version,
			})
		}
	} else {
		c.AbortWithStatus(401)
	}
}
