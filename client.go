package main
import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"strings"
	"encoding/json"
)

type Message struct {
	Type string
	Body string
}

func main() {
	msg, _ := json.Marshal(Message{
		Type: "Directory",
		Body: "/foo/bar.txt",
	})
    url := "http://127.0.0.1:8080"

	fmt.Println("Connecting to :8080 ...")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response Status:", resp.Status)
    fmt.Println("Response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", strings.Split(string(body), "\n"))
	// [http://127.0.0.1:8090/1 http://127.0.0.1:8090/2 http://127.0.0.1:8090/3 ]
}
