package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func CheckFileSize(data *commitData) (bool, error) {
	if data.GlobalWeb != 0 {
		fileType := "web"
		region := "global"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		fileSize, err := getFileSize(filePath)
		if err != nil {
			return false , err
		}
		if data.GlobalWeb != fileSize {
			return false, err
		}
		data.OldGlobalWeb = filePath
	}
	if data.GlobalProxy != 0 {
		fileType := "proxy"
		region := "global"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		fileSize, err := getFileSize(filePath)
		if err != nil {
			return false , err
		}
		if data.GlobalProxy != fileSize {
			return false, err
		}
		data.OldGlobalProxy = filePath
	}
	if data.GlobalDNS != 0 {
		fileType := "dns"
		region := "global"
		filename := "dns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		fileSize, err := getFileSize(filePath)
		if err != nil {
			return false , err
		}
		if data.GlobalDNS != fileSize {
			return false, err
		}
		data.OldGlobalDNS = filePath
	}
	if data.CNWeb != 0 {
		fileType := "web"
		region := "cn"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		fileSize, err := getFileSize(filePath)
		if err != nil {
			return false , err
		}
		if data.CNWeb != fileSize {
			return false, err
		}
		data.OldCNWeb = filePath
	}
	if data.CNProxy != 0 {
		fileType := "proxy"
		region := "cn"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		fileSize, err := getFileSize(filePath)
		if err != nil {
			return false , err
		}
		if data.CNProxy != fileSize {
			return false, err
		}
		data.OldCNProxy = filePath
	}
	return true, errors.New("something")
}

func getFileSize(filePath string) (int64, error){
	fi, err := os.Stat(filePath)
	check(err)
	if err != nil {
		return 0 , err
	}
	return fi.Size(), err
}

func StoreFile(data commitData) {
	version := data.Version
	if data.OldGlobalWeb != "" {
		newLocation := "data/web/global/"+version
		os.MkdirAll("data/web/global/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_web = version
	}
	if data.OldGlobalProxy != "" {
		newLocation := "data/proxy/global/"+version
		os.MkdirAll("data/proxy/global/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_web = version
	}
	if data.OldGlobalDNS != "" {
		newLocation := "data/dns/global/"+version
		os.MkdirAll("data/dns/global/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalWeb, newLocation+"/dns.tar.gz")
		check(err)
		serial_global_web = version
	}
	if data.OldCNWeb != "" {
		newLocation := "data/web/cn/"+version
		os.MkdirAll("data/web/cn/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_web = version
	}
	if data.OldCNProxy != "" {
		newLocation := "data/proxy/cn/"+version
		os.MkdirAll("data/proxy/cn/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_web = version
	}
}

func readSerial(target string) string {
	err := godotenv.Load("serialNumber")
	if err != nil {
		log.Fatal("Error loading serialNumber file")
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

func GetFile(versionType string, region string, version string) []byte {
	dir := fmt.Sprintf("data/%v/%v/%v/dns.tar.gz", versionType, region, version)

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println(info)
		log.Warn(info)
	}

	dat, err := ioutil.ReadFile(dir)
	check(err)
	return dat
}

func Md5sum(region string, versionType string) string {
	return "8a17a50573453a463ac8971079fafefa"
}

func check(e error) {
	if e != nil {
		fmt.Println("error:", e)
	}
}