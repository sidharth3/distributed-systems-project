package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	ip      = "127.0.0.1"
	TIMEOUT = 5
)

func IP() string {
	return ip + ":" + os.Args[1]
}

func StorageDir() string {
	return fmt.Sprintf("files_%v", os.Args[1])
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

func MasterIP() string {
	return os.Args[2]
}
