package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

var serial_web, serial_proxy, serial_dns string

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	token := os.Getenv("ACCEPTABLE_TOKEN")
	originToken := os.Getenv("ORIGIN_TOKEN")
	tempDir := os.Getenv("TEMP_DIR")
	dataDir := os.Getenv("DATA_DIR")

	serial_web = readSerial("web")
	serial_proxy = readSerial("proxy")
	serial_dns = readSerial("dns")

	const (
		layoutTimestamp = "2006-01-02 15:04:05"
	)

	//fmt.Println(time.Now().Unix())
	//fmt.Println(time.Now().UnixNano())

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

	})

	r.GET("/version/:type", GetSerial)
	r.GET("/version/:type/:version", Download)

	r.Run(":"+port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
