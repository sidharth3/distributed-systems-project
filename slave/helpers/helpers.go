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
