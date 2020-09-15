package cleaner

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func StartCleaner(duration time.Duration) {
	for true {
		cleanExpireData()
		time.Sleep(duration)
	}
}

func cleanExpireData() {
	dataType := []string{"dns", "proxy/global", "proxy/cn", "web/global", "web/cn", "all/global"}
	fmt.Println(fmt.Sprintf("%v - Clean job start", time.Now().Format("2006-01-02- 15:04:05")))
	for _, dt := range dataType {
		folders, err := ioutil.ReadDir(fmt.Sprintf("data/%v", dt))
		if err != nil {
			panic(err)
		} else {
			nowDate := time.Now().AddDate(0, 0, -1)
			atLeastOne := false
			removeList := []string{}
			for i := len(folders) - 1; i >= 0; i-- {
				folder := folders[i]
				folderDate, _ := time.Parse("20060102150405", folder.Name())
				//fmt.Println(folder.Name())
				if nowDate.After(folderDate) {
					removeList = append(removeList, folder.Name())
				} else {
					atLeastOne = true
				}
			}
			if atLeastOne == false && len(removeList) > 0 {
				removeList = append(removeList[1:])
			}
			for _, folder := range removeList {
				err := os.RemoveAll(fmt.Sprintf("data/%v/%v", dt, folder))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
