package broker

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	apiName    = "Norman - Broker APIs"
	apiVersion = "v0.0.1"
)

func (b *Broker) apiVersion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"API Name":    apiName,
			"API Version": apiVersion,
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		})
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (b *Broker) runQuery(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// err := json.NewDecoder(r.Body).Decode(&request)
		// if err != nil {
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// }

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"result":    "nice job!",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
