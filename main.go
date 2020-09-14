package main

import (
	"LeaCenter/cleaner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var serial_cn_web, serial_global_web, md5_cn_web, md5_global_web,
	serial_cn_proxy, serial_global_proxy, md5_cn_proxy, md5_global_proxy,
	serial_global_dns, md5_global_dns, serial_global_all, md5_global_all string

var isOrigin bool

var tempDir, dataDir string

type commitData struct {
	Version string `json:"version"`
	WebVersion string `json:"webVersion"`
	ProxyVersion string `json:"proxyVersion"`
	DnsVersion string `json:"dnsVersion"`
	TmpPathGlobalWeb string
	TmpPathGlobalProxy string
	TmpPathGlobalDNS string
	TmpPathGlobalAll string
	TmpPathCNWeb string
	TmpPathCNProxy string
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
	serial_global_all = readSerial("GLOBAL_ALL")

	cleaner_interval := readSerial("CLEANER_INTERVAL")

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

	//go cleaner.StartCleaner(24*time.Hour)
	if cleaner_interval != "" {
		cleaner_interval, _ := time.ParseDuration(cleaner_interval)
		go cleaner.StartCleaner(cleaner_interval)
	}

	r.Run(":"+port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
