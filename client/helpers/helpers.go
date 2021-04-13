package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
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
