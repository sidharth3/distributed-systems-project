package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
)

var clientPort string

func StorageDir() string {
	return fmt.Sprintf("files_%v", clientPort)
}

func MockClient() {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", b)
	}))
	clientPort = getPortNumber(ts.URL)
	fmt.Println("Successfully started client at", ts.URL)
}

func getPortNumber(clientURL string) string {
	u, err := url.Parse(clientURL)
	if err != nil {
		panic(err)
	}
	_, port, _ := net.SplitHostPort(u.Host)
	return port
}

func OpenFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
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
	h.Write(byteArr[:numBytes])
	hashValue := hex.EncodeToString(h.Sum(nil))

	return hashValue
}
