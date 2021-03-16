package SlaveOps

import (
	"net/http"
	"fmt"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.Form["file"][0]
	http.ServeFile(w, r, fmt.Sprintf("files/%v", filename))
	// http.ServeFile(w, r, "files/test_file.txt")
}
