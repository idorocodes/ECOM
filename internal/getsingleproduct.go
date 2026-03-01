package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Id string `json:"id"`
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Request recieved by /getaproduct")

		decoder := json.NewDecoder(r.Body)

		var reqBody Request
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
			return
		}

		if len(reqBody.Id) == 0 {
			http.Error(w, "Id not supplied", http.StatusBadRequest)
			return
		}

		response, error := GetSingleProduct(reqBody.Id)
		if error != nil {

			http.Error(w, fmt.Sprintf("%v", error), http.StatusInternalServerError)
			return
		} else {

			data := response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
}
