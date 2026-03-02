package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteRequest struct {
	Id string `json:"id"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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
		fmt.Println("Request recieved by /deleteProduct")

		decoder := json.NewDecoder(r.Body)

		var reqBody DeleteRequest

		err := decoder.Decode(&reqBody)

		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
		}

		if len(reqBody.Id) == 0 {
			http.Error(w, "Id not supplied!", http.StatusBadRequest)
		}

		response, error := DeleteSingleProduct(reqBody.Id)

		if error != nil || !response {
			http.Error(w, fmt.Sprintf("%v", error), http.StatusInternalServerError)
		}

		if response == true {
			data := DeleteResponse{
				Message: "Product deleted from the database",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}
