package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

var serial_cn_web, serial_global_web, md5_cn_web, md5_global_web,
	serial_cn_proxy, serial_global_proxy, md5_cn_proxy, md5_global_proxy,
	serial_global_dns, md5_global_dns string

var isOrigin bool

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	token := os.Getenv("ACCEPTABLE_TOKEN")
	originToken := os.Getenv("ORIGIN_TOKEN")
	tempDir := os.Getenv("TEMP_DIR")
	dataDir := os.Getenv("DATA_DIR")

	serial_cn_web = readSerial("GLOBAL_WEB")
	serial_global_web = readSerial("CN_WEB")
	serial_cn_proxy = readSerial("GLOBAL_PROXY")
	serial_global_proxy = readSerial("CN_PROXY")
	serial_global_dns = readSerial("GLOBAL_DNS")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if port == "" {
		log.Fatal("PORT is empty in .env file")
	} else if token == "" {
		log.Fatal("ACCEPTABLE_TOKEN is empty in .env file")
	} else if tempDir == "" {
		log.Fatal("TEMP_DIR is empty in .env file")
	} else if dataDir == "" {
		log.Fatal("DATA_DIR is empty in .env file")
	}

	r := gin.Default()
	r.Use(AuthToken(token, originToken))

	r.GET("/ping", Pong)

	r.POST("/version/commit", func(c *gin.Context) {
		type jsonData struct {
			Version string `json:"version"`
			GlobalWeb int `json:"webGlobal"`
			GlobalProxy int `json:"ProxyGlobal"`
			GlobalDNS int `json:"dnsGlobal"`
			CNWeb int `json:"webCN"`
			CNProxy int `json:"proxyCN"`


		}
		var data jsonData
		c.BindJSON(&data)
		fmt.Println(data)
		fmt.Println(data.Version)
		fmt.Println(data.GlobalWeb)
		fmt.Println(data.GlobalProxy)
		fmt.Println(data.GlobalDNS)
		fmt.Println(data.CNWeb)
		fmt.Println(data.CNProxy)
		c.JSON(200, gin.H{
			"version": data.Version,
			"webGlobal": data.GlobalWeb,
			"proxyGlobal": data.GlobalProxy,
			"dnsGlobal": data.GlobalDNS,
			"webCN": data.CNWeb,
			"proxyCN": data.CNProxy,
		})
	})

	r.GET("/version/:region/:type", GetSerial)
	r.GET("/version/:region/:type/:version", Download)

	r.Run(":"+port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
