package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

type CreateAcc struct {
	FirstName       string `json:"firstname"`
	SecondName      string `json:"secondname"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
}

type CreateAccResponse struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Message    string `jsosn:"message"`
	Token      string `json:"token"`
	Success    bool   `json:"success"`
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

		fmt.Println("Request sent to create account")
		decoder := json.NewDecoder(r.Body)

		var reqBody CreateAcc
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad Body", http.StatusBadRequest)
		}

		if reqBody.ConfirmPassword != reqBody.Password {
			error := "Password does not match"
			http.Error(w, error, http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(reqBody.Email)
		if err != nil {
			errors := "Email is invalid"
			http.Error(w, errors, http.StatusBadRequest)
			return
		}

		if len(reqBody.FirstName) == 0 || len(reqBody.SecondName) == 0 || len(reqBody.Username) == 0 {
			errors := "Name is invalid"
			http.Error(w, errors, http.StatusBadRequest)
			return
		}

		
		token, err := CreateToken(reqBody.Username)

		data := CreateAccResponse{
			Message:    "Account Created!",
			FirstName:  reqBody.FirstName,
			SecondName: reqBody.SecondName,
			Email:      reqBody.Email,
			Username:   reqBody.Username,
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
