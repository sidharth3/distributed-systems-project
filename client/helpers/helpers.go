package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"regexp"
	"strings"
)

func OpenFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func HashFileContent(f *os.File) string {
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	byteArr := make([]byte, fi.Size())
	numBytes, err := f.Read(byteArr)
	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	_, err = h.Write(byteArr[:numBytes])
	if err != nil {
		log.Fatal(err)
	}
	hashValue := hex.EncodeToString(h.Sum(nil))

	return hashValue
}

func SanitizeInput(path string) string {
	sanitizedPath := ""
	pathList := strings.Split(path, "/")
	check, err := regexp.MatchString(`.\..`, pathList[len(pathList)-1])
	if !check {
		log.Fatal("Filename invalid")
	}
	if err != nil {
		log.Fatal(err)
	}
	if check {

	}
	for i, path := range pathList {
		if i != len(pathList)-1 {
			path = strings.ReplaceAll(path, ".", "")
		}
		if strings.Trim(path, " ") != "" {
			sanitizedPath += "/" + path
		}
	}
	return sanitizedPath
}
