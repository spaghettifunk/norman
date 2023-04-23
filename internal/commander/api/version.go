package api

import "net/http"

func Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Norman Commander API Version 0.0.1"))
	w.WriteHeader(http.StatusOK)
}
