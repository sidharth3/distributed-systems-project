package handlers

import (
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func GarbageCollectorHandler(w http.ResponseWriter, r *http.Request){
	filesBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	files := make( map[string]bool)
	err = json.Unmarshal(filesBytes, &files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RECIEVED:",files)
	//check if in, else delete
	dirs := helpers.ListDir()
	for dir := range dirs{
		if !files[dir]{
			fmt.Println("File",dir,"is no longer referenced in master, and will be deleted.")
			helpers.DeleteFile(dir)
		}
	}

}
