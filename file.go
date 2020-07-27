package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)


func storeFile(serialNumber string) {

}

func readSerial(target string) string {
	err := godotenv.Load("serialNumber")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serialNumber := os.Getenv(target)
	if serialNumber != "" {
		log.Info(target + " " + serialNumber)
	} else {
		fmt.Println(serialNumber)
		log.Warning("empty "+target+" value")
	}
	return serialNumber
}

func GetFile(versionType string, version string) []byte {
	dir := fmt.Sprintf("data/%v/%v/leadns.tar.gz", versionType, version)

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println(info)
		log.Warn(info)
	}

	dat, err := ioutil.ReadFile(dir)
	check(err)
	return dat
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}