package handlers

import (
	"ds-proj/master/structs"
	"encoding/json"
	"log"
	"net/http"
)

func HandleFile(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		w.Header().Set("Content-Type", "application/json")
		ipArr := make([]string, 0)
		for _, slave := range m.DirectoryTable[filename] {
			ipArr = append(ipArr, slave.IP)
		}
		data, err := json.Marshal(ipArr)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}
