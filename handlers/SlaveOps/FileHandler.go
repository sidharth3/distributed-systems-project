package SlaveOps

import (
	"fmt"
	"net/http"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.Form["file"][0]
	http.ServeFile(w, r, fmt.Sprintf("../../files/%v", filename))
}
