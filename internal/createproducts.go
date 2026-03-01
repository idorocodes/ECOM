package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	Name            string `json:"name"`
	Price           int    `json:"price"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	DefaultCurrency string `json:"defaultcurrency"`
}

type GetProduct struct {
	Id 				string `json:"id"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	DefaultCurrency string `json:"defaultcurrency"`
}

type CreateProductResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		fmt.Println("Request recieved by /createProduct")
		decode := json.NewDecoder(r.Body)

		var reqBody Product

		err := decode.Decode(&reqBody)

		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
			return
		}

		if len(reqBody.Name) == 0 || len(reqBody.Description) == 0 || len(reqBody.Category) == 0 || len(reqBody.DefaultCurrency) == 0 {
			http.Error(w, "Check body", http.StatusBadRequest)
			return
		}
		reponse, error := CreateAProduct(reqBody.Name, reqBody.Price, reqBody.Description, reqBody.Category, reqBody.DefaultCurrency)

		if error != nil {

			http.Error(w, fmt.Sprintf("%v", error), http.StatusInternalServerError)
			return
		} else {

			data := CreateProductResponse{
				Message: reponse,
				Success: true,
				Code:    http.StatusCreated,
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
