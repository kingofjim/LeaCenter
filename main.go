package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

var serial_cn_web, serial_global_web, md5_cn_web, md5_global_web,
	serial_cn_proxy, serial_global_proxy, md5_cn_proxy, md5_global_proxy,
	serial_global_dns, md5_global_dns string

var isOrigin bool

var tempDir, dataDir string

type commitData struct {
	Version string `json:"version"`
	GlobalWeb int64 `json:"webGlobal"`
	GlobalProxy int64 `json:"ProxyGlobal"`
	GlobalDNS int64 `json:"dnsGlobal"`
	CNWeb int64 `json:"webCN"`
	CNProxy int64 `json:"proxyCN"`
	OldGlobalWeb string
	OldGlobalProxy string
	OldGlobalDNS string
	OldCNWeb string
	OldCNProxy string
}

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	token := os.Getenv("ACCEPTABLE_TOKEN")
	originToken := os.Getenv("ORIGIN_TOKEN")
	tempDir = os.Getenv("TEMP_DIR")
	dataDir = os.Getenv("DATA_DIR")

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
	} else if originToken == "" {
		log.Fatal("ORIGIN_TOKEN is empty in .env file")
	}

	r := gin.Default()
	r.Use(AuthToken(token, originToken))

	r.GET("/ping", Pong)
	r.POST("/version/commit", Commit)
	r.GET("/version/:type/:region", GetSerial)
	r.GET("/version/:type/:region/:version", Download)

	r.Run(":"+port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
