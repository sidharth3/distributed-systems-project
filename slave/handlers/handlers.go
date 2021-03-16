package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.Form["file"][0]
	http.ServeFile(w, r, filepath.Join(fmt.Sprintf("files_%v", os.Args[1]), filename))
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "i is ok")
}
