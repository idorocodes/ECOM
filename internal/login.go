package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

type LoginAcc struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAccResponse struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Message    string `jsosn:"message"`
	Token      string `json:"token"`
	Role       string `json:"role"`
	Success    bool   `json:"success"`
}

func LoginAccount(w http.ResponseWriter, r *http.Request) {
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

		fmt.Println("Request recieved by /loginAccount")
		decoder := json.NewDecoder(r.Body)

		var reqBody LoginAcc
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
		}

		_, err = mail.ParseAddress(reqBody.Email)
		if err != nil {
			errors := "Email is invalid"
			http.Error(w, errors, http.StatusBadRequest)
			return
		}

		hashedPassword := HashPassowrd(reqBody.Password)

		response, error := LoginUser(hashedPassword, reqBody.Email)

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
