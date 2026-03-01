package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateRequest struct {
	Id       string `json:"id"`
	NewName  string `json:"newname"`
	NewPrice int    `json:"newprice"`
}

type UpdateResponse struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	DefaultCurrency string `json:"defaultcurrency"`
	Message         string `json:"message"`
}

func UpdateAProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
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
		fmt.Println("Request recieved by /updateaproduct")

		decoder := json.NewDecoder(r.Body)

		var reqBody UpdateRequest
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
			return
		}

		if len(reqBody.Id) == 0 || len(reqBody.NewName) == 0 || reqBody.NewPrice == 0 {
			http.Error(w, "Worng body data!", http.StatusBadRequest)
			return
		}

		response, error := UpdateSingleProduct(reqBody.Id, reqBody.NewName, reqBody.NewPrice)
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
