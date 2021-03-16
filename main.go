package main

import (
	"ds-proj/handlers/SlaveOps"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const PING_INTERVAL = 3 // in s

// Declaring node structures
type Node struct {
	Port int
}

type Master struct {
	Node
	Slaves         map[*Slave]Status // Slave instance => Status
	DirectoryTable map[string]*Slave // "/foo/bar.txt" => Slave instance
}

type Slave struct {
	Node
	ID int
	*Master
	Files map[string][]byte // "/foo/bar.txt" => byte array
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
	DEAD        Status = "Cannot be pinged"
)

// Master Functions
func ConstructMaster(port int) *Master {
	return &Master{Node{port}, make(map[*Slave]Status), make(map[string]*Slave)}
}

func (m *Master) HttpServerStart() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", m.Port), nil))
}

func (m *Master) PingSlave(s *Slave) {
	url := fmt.Sprintf("http://localhost:%d/ping", s.Port)
	fmt.Println("Master is pinging Slave", s.ID, "periodically at", url)

	ticker := time.NewTicker(PING_INTERVAL * time.Second)
	go func() {
		for {
			<-ticker.C
			// Send the GET request
			res, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			}

			// Read the response
			reply, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(reply))
		}
	}()
}

// Slave Functions
func (s *Slave) HttpServerStart() {
	http.HandleFunc("/ping", SlaveOps.HeartbeatHandler)
	http.HandleFunc("/file", SlaveOps.DownloadFile)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))
}

func ConstructSlave(port int, id int, master *Master) *Slave {
	return &Slave{Node{port}, id, master, make(map[string][]byte)}
}

func main() {
	const START_PORT int = 8080

	// Set up 1 master and 1 slave for now
	m := ConstructMaster(START_PORT)
	s := ConstructSlave(START_PORT+1, 1, m)
	m.Slaves[s] = UNDERLOADED

	// Start all the nodes
	go s.HttpServerStart()
	//go m.HttpServerStart()

	// Test heartbeat message
	// time.Sleep(1000)
	// for slave := range m.Slaves {
	// 	m.PingSlave(slave)
	// }

	// Prevents the program from terminating
	var input string
	fmt.Scanln(&input)
}
