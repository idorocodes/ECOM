package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

type CreateAcc struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Password   string `json:"password"`
}

type CreateAccResponse struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `jons:"role"`
	Message    string `jsosn:"message"`
	Token      string `json:"token"`
	Success    bool   `json:"success"`
}

type AlreadyCreatedResponse struct {
	Message string `json:"message"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
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

		fmt.Println("Request recieved by /createAccount")
		decoder := json.NewDecoder(r.Body)

		var reqBody CreateAcc
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

		if len(reqBody.FirstName) == 0 || len(reqBody.SecondName) == 0 || len(reqBody.Username) == 0 || len(reqBody.Role) == 0 {
			errors := "Name  is invalid "
			http.Error(w, errors, http.StatusBadRequest)
			return
		}

		if len(reqBody.Role) == 0 || reqBody.Role != "customer" && reqBody.Role != "admin" {
			errors := " Role is invalid "
			http.Error(w, errors, http.StatusBadRequest)
			return
		}

		token, err := CreateToken(reqBody.Username, reqBody.Role)

		hashedPassword := HashPassowrd(reqBody.Password)

		reponse, error := CreateUser(reqBody.Username, reqBody.FirstName, reqBody.SecondName, hashedPassword, reqBody.Email, reqBody.Role)

		if fmt.Sprintf("%v", error) == "User already exists in the db" {
			data := AlreadyCreatedResponse{
				Message: "User already exist in the database, please login",
			}
			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			return
		} else if error != nil {

			http.Error(w, fmt.Sprintf("%v", error), http.StatusInternalServerError)
			return
		} else {
			data := CreateAccResponse{
				Message:    reponse,
				FirstName:  reqBody.FirstName,
				SecondName: reqBody.SecondName,
				Email:      reqBody.Email,
				Username:   reqBody.Username,
				Role:       reqBody.Role,
				Success:    true,
				Token:      token,
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
