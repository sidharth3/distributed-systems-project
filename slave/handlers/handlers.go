package handlers

import (
	"ds-proj/slave/helpers"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.Form["file"][0]
	http.ServeFile(w, r, filepath.Join(helpers.StorageDir(), filename))
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(helpers.ListDir())
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}
