package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)




func GetProducts(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Request recieved by /getProducts")
		
		reponse, error := GetAllProducts()

		if error != nil {

			http.Error(w, fmt.Sprintf("%v", error), http.StatusInternalServerError)
			return
		} else {

			data := reponse
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
}
