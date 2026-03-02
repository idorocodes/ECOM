package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderRequest struct {
	ProductId string `json:"productid"`
	Address   string `json:"address"`
}

type OrderResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Sucess  bool   `json:"success"`
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	type contextKey string
	const userIdKey contextKey = "userId"

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
		fmt.Println("Request recieved by /createOrder")

		decoder := json.NewDecoder(r.Body)

		var reqBody OrderRequest

		err := decoder.Decode(&reqBody)

		if err != nil {
			http.Error(w, "Json Error", http.StatusInternalServerError)

			return
		}

		if len(reqBody.Address) == 0 || len(reqBody.ProductId) == 0 {
			http.Error(w, "Bad body, some  fields are not supplied", http.StatusBadRequest)
			return
		}

		userId, ok := r.Context().Value(ctxId{}).(string)
	

		if !ok {
			http.Error(w, "missing user id", http.StatusInternalServerError)
			return

		}

		response, code, err := CreateAnOrder(reqBody.ProductId, userId, reqBody.Address)

		if fmt.Sprintf("%v", err) == "No products found" {
			http.Error(w, "Product not found in the database!", http.StatusNotFound)
			return
		}
		if err != nil {

			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		} else {

			data := OrderResponse{
				Message: response,
				Code:    code,
				Sucess:  true,
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
