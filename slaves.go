package main
import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Message struct {
	Type string
	Body string
}

func conn(w http.ResponseWriter, req *http.Request){
	fmt.Println("Handling a /1 connection ...")

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	msg := buf.String()
	fmt.Println("Msg: ", msg)

	if msg == "Alive"{
		w.Write([]byte("Yes"))
	}
}

func register(IP string){
	msg, _ := json.Marshal(Message{
		Type: "Register",
		Body: IP,
	})
    url := "http://127.0.0.1:8080"

	fmt.Println("Connecting to :8080 ...")
    var location = []byte(msg)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(location))

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

}

func main(){
	fmt.Println("Starting slaves ...")
	register("http://127.0.0.1:8090/1")
	http.HandleFunc("/1",conn)
	http.ListenAndServe("127.0.0.1:8090", nil)
	fmt.Println("Listening on :8090 ...")
}