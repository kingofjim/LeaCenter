package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
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
		filename := "dns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+filename
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
	if data.CNProxy == 0 && data.CNWeb == 0 && data.GlobalDNS == 0 && data.GlobalWeb == 0 && data.GlobalProxy == 0 {
		return false, errors.New("All parameter is 0")
	}
	return true, nil
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
		err := os.Rename(data.OldGlobalProxy, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_proxy = version
	}
	if data.OldGlobalDNS != "" {
		newLocation := "data/dns/"+version
		os.MkdirAll("data/dns/"+version, os.ModePerm)
		err := os.Rename(data.OldGlobalDNS, newLocation+"/dns.tar.gz")
		check(err)
		serial_global_dns = version
	}
	if data.OldCNWeb != "" {
		newLocation := "data/web/cn/"+version
		os.MkdirAll("data/web/cn/"+version, os.ModePerm)
		err := os.Rename(data.OldCNWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_cn_web = version
	}
	if data.OldCNProxy != "" {
		newLocation := "data/proxy/cn/"+version
		os.MkdirAll("data/proxy/cn/"+version, os.ModePerm)
		err := os.Rename(data.OldCNProxy, newLocation+"/leadns.tar.gz")
		check(err)
		serial_cn_proxy = version
	}
	CleanTmp()
	go UpdateSerialFile()
}

func readSerial(target string) string {
	err := godotenv.Load("serialNumber")
	if err != nil {
		log.Warning("Error loading serialNumber file")
	}
	serialNumber := os.Getenv(target)
	if serialNumber != "" {
		log.Info(target + " " + serialNumber)
	} else {
		log.Warning("empty "+target+" value")
	}
	return serialNumber
}

func GetFile(versionType string, region string, version string) []byte {
	var filename, dir string
	if versionType == "dns" {
		filename = "dns.tar.gz"
		dir = fmt.Sprintf("data/%v/%v/%v", versionType, version, filename)

	} else {
		filename = "leadns.tar.gz"
		dir = fmt.Sprintf("data/%v/%v/%v/%v", versionType, region, version, filename)
	}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Warn(info)
	}

	dat, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}
	return dat
}

func Md5sum(versionType string, region string, version string) string {
	var filename, path string
	if versionType == "dns" {
		filename = "dns.tar.gz"
		path = fmt.Sprintf("data/%v/%v/%v", versionType, version, filename)
	} else {
		filename = "leadns.tar.gz"
		path = fmt.Sprintf("data/%v/%v/%v/%v", versionType, region, version, filename)

	}

	f, err := os.Open(path)
	if err != nil {
		log.Warning(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Warning(err)
	}

	md5 := fmt.Sprintf("%x", h.Sum(nil))

	return md5
}

func CleanTmp() {
	err := os.RemoveAll(tempDir+"/web")
	check(err)
	err = os.RemoveAll(tempDir+"/proxy")
	check(err)
	err = os.RemoveAll(tempDir+"/dns")
	check(err)
}

func UpdateSerialFile() {
	text := fmt.Sprintf("GLOBAL_WEB=%v\n" +
		"GLOBAL_PROXY=%v\n" +
		"GLOBAL_DNS=%v\n" +
		"CN_WEB=%v\n" +
		"CN_PROXY=%v", serial_global_web, serial_global_proxy, serial_global_dns, serial_cn_web, serial_cn_proxy)
	env, err := godotenv.Unmarshal(text)
	check(err)
	err = godotenv.Write(env, "./serialNumber")
	check(err)
}

func check(e error) {
	if e != nil {
		log.Panicln(e)
	}
}