package handlers

import (
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.Form["file"][0]
	http.ServeFile(w, r, filepath.Join(helpers.StorageDir(), filename))
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	file, fileHeader, err := r.FormFile("filename")
	// uid, err := r.FormFile("uid")
	uid := fmt.Sprint(r.Form["uid"])
	fmt.Println(fileHeader.Filename)
	fmt.Println(uid)
	data := url.Values{"filename": {fileHeader.Filename}, "uid": {fmt.Sprint(uid)}}
	ForceUpdateMaster(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Create a new file in the uploads directory
	dst, err := os.Create(path.Join(helpers.StorageDir(), fileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded File: %v\n", fileHeader.Filename)
	fmt.Printf("File Size: %v\n", fileHeader.Size)

	// TODO: need to inform master first?

	// return that we have successfully uploaded our file
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func ForceUpdateMaster(data url.Values) {
	//Questions - does forceUpdate need to send directory also or new upload information only?
	master_URL := "http://127.0.0.1:8080/update"
	res, err := http.PostForm(master_URL, data)
	fmt.Println(res.StatusCode)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("File upload has failed.")
	}
}

func HandleReplica(w http.ResponseWriter, r *http.Request) {
	toGetByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	toGet := make(map[string]string)
	err = json.Unmarshal(toGetByte, &toGet)

	// toGet => {fileHash: ip1, fileHash, ip2}
	for f, ip := range toGet {
		res, err := http.Get("http://" + ip + "/file?file=" + f)
		if err != nil {
			log.Fatal(err)
		}

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

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	w.WriteHeader(http.StatusOK)
	files := make(map[string]bool)
	err = json.Unmarshal(filesBytes, &files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RECIEVED:", files)
	//check if in, else delete
	dirs := helpers.ListDir()
	for dir := range dirs {
		if !files[dir] {
			fmt.Println("File", dir, "is no longer referenced in master, and will be deleted.")
			helpers.DeleteFile(dir)
		}
	}

}
