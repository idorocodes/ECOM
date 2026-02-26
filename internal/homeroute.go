package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
}



func HomePath(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		data := Response{
			Message: "Route does not exist for this method",
			Code:    http.StatusMethodNotAllowed,
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {

		fmt.Println("Request received on path / ")
		data := Response{
			Message: "Welcome to ECOM",
			Code:    http.StatusOK,
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
