package SlaveOps

import (
	"fmt"
	"net/http"
)

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "i is ok")
}
