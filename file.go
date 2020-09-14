package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

func CheckFileExist(data *commitData) (bool, error) {
	if data.WebVersion != "0" {
		fileType := "web"
		region := "global"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathGlobalWeb = filePath
		} else {
			return false, err
		}

		region = "cn"
		filename = "leadns.tar.gz"
		filePath = tempDir+"/"+fileType+"/"+region+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathCNWeb = filePath
		} else {
			return false, err
		}
	}
	if data.ProxyVersion != "0" {
		fileType := "proxy"
		region := "global"
		filename := "leadns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+region+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathGlobalProxy = filePath
		} else {
			return false, err
		}
		region = "cn"
		filename = "leadns.tar.gz"
		filePath = tempDir+"/"+fileType+"/"+region+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathCNProxy = filePath
		} else {
			return false, err
		}
	}
	if data.DnsVersion != "0" {
		fileType := "dns"
		filename := "dns.tar.gz"
		filePath := tempDir+"/"+fileType+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathGlobalDNS = filePath
		} else {
			return false, err
		}
	}

	if data.Version != "0" {
		filename := "compressfile.tar.gz"
		filePath := tempDir+"/"+filename
		if _, err := os.Stat(filePath); err == nil {
			data.TmpPathGlobalAll = filePath
		} else {
			return false, err
		}
	}
	return true, nil
}

func StoreFile(data commitData) {
	version := data.Version
	if data.WebVersion != "0" && data.TmpPathGlobalWeb != "" {
		newLocation := "data/web/global/"+version
		os.MkdirAll("data/web/global/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathGlobalWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_web = version
	}
	if data.ProxyVersion != "0" && data.TmpPathGlobalProxy != "" {
		newLocation := "data/proxy/global/"+version
		os.MkdirAll("data/proxy/global/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathGlobalProxy, newLocation+"/leadns.tar.gz")
		check(err)
		serial_global_proxy = version
	}
	if data.DnsVersion != "0" && data.TmpPathGlobalDNS != "" {
		newLocation := "data/dns/"+version
		os.MkdirAll("data/dns/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathGlobalDNS, newLocation+"/dns.tar.gz")
		check(err)
		serial_global_dns = version
	}
	if data.WebVersion != "0" && data.TmpPathCNWeb != "" {
		newLocation := "data/web/cn/"+version
		os.MkdirAll("data/web/cn/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathCNWeb, newLocation+"/leadns.tar.gz")
		check(err)
		serial_cn_web = version
	}
	if data.ProxyVersion != "0" && data.TmpPathCNProxy != "" {
		newLocation := "data/proxy/cn/"+version
		os.MkdirAll("data/proxy/cn/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathCNProxy, newLocation+"/leadns.tar.gz")
		check(err)
		serial_cn_proxy = version
	}
	if data.TmpPathGlobalAll != "" {
		newLocation := "data/all/global/"+version
		os.MkdirAll("data/all/global/"+version, os.ModePerm)
		err := os.Rename(data.TmpPathGlobalAll, newLocation+"/compressfile.tar.gz")
		check(err)
		serial_global_all = version
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
	if versionType == "all" {
		filename = "compressfile.tar.gz"
		dir = fmt.Sprintf("data/%v/%v/%v/%v", versionType, region, version, filename)
	} else if versionType == "dns" {
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
	if versionType == "all" {
		filename = "compressfile.tar.gz"
		path = fmt.Sprintf("data/%v/%v/%v/%v", versionType, region, version, filename)
	} else if versionType == "dns" {
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
	err = os.RemoveAll(tempDir+"/compressfile.tar.gz")
	check(err)
}

func UpdateSerialFile() {
	text := fmt.Sprintf("GLOBAL_WEB=%v\n" +
		"GLOBAL_PROXY=%v\n" +
		"GLOBAL_DNS=%v\n" +
		"CN_WEB=%v\n" +
		"CN_PROXY=%v\n" +
		"GLOBAL_ALL=%v", serial_global_web, serial_global_proxy, serial_global_dns, serial_cn_web, serial_cn_proxy, serial_global_all)
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

func check_error_500(c *gin.Context, e error) bool{
	if e != nil {
		log.Panicln(e)
		status := 500
		fmt.Println(e)
		c.JSON(status, gin.H{
			"error": e.Error(),
		})
		return true
	}
	return false
}