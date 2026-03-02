package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BadResponse struct {
	Message string `json:"message"`
	Code    int    `json:"int"`
	Success bool   `json:"success"`
}

type AllUsers struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type AllOrders struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Address   string `json:"address"`
	UserId    string `json:"user_id"`
}

type MetricResponse struct {
	Users    []AllUsers  `json:"allusers"`
	Products []Product   `json:"product"`
	Orders   []AllOrders `json:"orders"`
	Success  bool        `json:"success"`
}

func Metrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		data := BadResponse{
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
		fmt.Println("Request recieved by /get Metrics")

		allUsers, err := GetAllUsers()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allProducts, err := GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allOrders, err := GetAllOrders()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := MetricResponse{
			Users:    allUsers,
			Products: allProducts,
			Orders:   allOrders,
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
