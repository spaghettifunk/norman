package api

import (
	"encoding/json"
	"net/http"
)

type versionResponse struct {
	Message string `json:"message"`
}

func Version() http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)

		vr := &versionResponse{Message: "Norman Commander API Version 0.0.1"}
		b, err := json.Marshal(vr)
		if err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			rw.Write([]byte("Could not unmarshall versionResponse object"))
			return
		}
		rw.Write(b)
	}
	return http.HandlerFunc(fn)
}
