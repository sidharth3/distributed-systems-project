package handlers

import (
	"ds-proj/slave/config"
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename := r.Form["file"][0]
	http.ServeFile(w, r, filepath.Join(helpers.StorageDir(), filename))
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("filename")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the uploads directory
	dst, err := os.Create(path.Join(helpers.StorageDir(), fileHeader.Filename))
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleReplica(w http.ResponseWriter, r *http.Request) {
	toGetByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	toGet := make(map[string]string)
	err = json.Unmarshal(toGetByte, &toGet)
	if err != nil {
		log.Fatal(err)
	}

	// toGet => {fileHash: ip1, fileHash, ip2}
	for f, ip := range toGet {
		client := &http.Client{
			Timeout: time.Second * config.TIMEOUT,
		}
		res, _ := client.Get("http://" + ip + "/file?file=" + f)

		// Create the file
		if res.StatusCode == http.StatusOK {
			out, err := os.Create(path.Join(helpers.StorageDir(), f))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// Write the body to file
			_, err = io.Copy(out, res.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func LoadHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(helpers.GetLoad())
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(helpers.ListDir())
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func GarbageCollectorHandler(w http.ResponseWriter, r *http.Request) {
	filesBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	files := make(map[string]bool)
	err = json.Unmarshal(filesBytes, &files)
	if err != nil {
		log.Fatal(err)
	}
	//Delete if in
	dirs := helpers.ListDir()
	for dir := range dirs {
		if files[dir] {
			fmt.Println("File", dir, "is no longer referenced in master, and will be deleted.")
			helpers.DeleteFile(dir)
		}
	}
}
