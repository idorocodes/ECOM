package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusResponse struct {
	Message string `jsosn:"message"`
	Success bool   `json:"success"`
}

func HealthStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		data := StatusResponse{
			Message: "Route does not exist for this method",
			Success: false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		fmt.Println("Request recieved by /healthStatus")
		data := StatusResponse{
			Message: "ECOM is working perfectly",
			Success: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
