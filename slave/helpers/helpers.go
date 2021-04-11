package helpers

import (
	"ds-proj/slave/config"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func IP() string {
	return config.BASEIP + ":" + os.Args[1]
}

func StorageDir() string {
	return fmt.Sprintf("files_%v", os.Args[1])
}

func GetLoad() int {
	files, err := ioutil.ReadDir(StorageDir())
	if err != nil {
		log.Fatal(err)
	}
	totalLoad := 0
	for _, file := range files {
		totalLoad += int(file.Size())
	}
	return totalLoad
}

func ListDir() map[string]bool {
	files, err := ioutil.ReadDir(StorageDir())
	if err != nil {
		log.Fatal(err)
	}
	filenames := make(map[string]bool)
	for _, file := range files {
		filenames[file.Name()] = true
	}
	return filenames
}

func DeleteFile(filename string) {
	e := os.Remove(StorageDir() + "/" + filename)
	if e != nil {
		log.Fatal(e)
	}
}

// func MasterIP() string {
// 	return os.Args[2]
// }

// MasterIP() returns a []str{} of master ips
func MasterIP() []string {
	iplist := make([]string, 0)

	for i := 0; i < len(os.Args)-2; i++ {
		if os.Args[i] != "" {
			iplist = append(iplist, os.Args[i+2])
		}
	}
	return iplist
}
