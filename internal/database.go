package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	SupaURL = "https://wampmtueakfcnvdvaxry.supabase.co/rest/v1"
	SupaKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndhbXBtdHVlYWtmY252ZHZheHJ5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzIxODM4MjgsImV4cCI6MjA4Nzc1OTgyOH0.EJFbo7C748DMu5fDl7WDm4qXliMkIQnrOb3reUL8n7k"
)

type User struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Message    string `jsosn:"message"`
}

func checkIfUserExists(email string) (bool, error) {

	url := fmt.Sprintf("%s/users?email=eq.%s&select=*", SupaURL, email)
	req, err := http.NewRequest("GET", url, nil)
	var result bool
	type User struct {
		FirstName       string `json:"firstname"`
		SecondName      string `json:"secondname"`
		Username        string `json:"username"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
	}

	if err != nil {
		return false, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return false, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		result = false

	} else {
		result = true
	}

	return result, nil
}

func CreateUser(username, firstname, secondname, password, email, role string) (string, error) {
	result, err := checkIfUserExists(email)

	if err != nil {
		return "", errors.New("Response error")
	}

	if result == true {
		return "", errors.New("User already exists in the db ")
	} else {
		url := SupaURL + "/users"
		var finalResponse string
		userData := map[string]string{
			"username":   username,
			"firstname":  firstname,
			"secondname": secondname,
			"password":   password,
			"email":      email,
			"role":       role,
		}

		data, err := json.Marshal(userData)

		if err != nil {
			return "", errors.New("Json error")
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			return "", errors.New("Request error")
		}

		req.Header.Set("apiKey", SupaKey)
		req.Header.Set("Authorization", "Bearer "+SupaKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Prefer", "return=minimal")

		client := &http.Client{}

		response, err := client.Do(req)

		if err != nil {
			return "", errors.New("Response error")
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(response.Body)
			return "", errors.New("Supabase Error" + string(body))
		}
		if response.StatusCode == http.StatusCreated {
			finalResponse = "User Created in the db"
		}

		return finalResponse, nil
	}
}

func LoginUser(password, email string) (LoginAccResponse, error) {
	url := fmt.Sprintf("%s/users?email=eq.%s&select=*", SupaURL, email)
	req, err := http.NewRequest("GET", url, nil)

	var loginResponse LoginAccResponse
	if err != nil {
		return LoginAccResponse{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return LoginAccResponse{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return LoginAccResponse{}, errors.New("User does not exist , please register")

	} else {

		dbUserPassword := results[0]["password"].(string)

		if password != dbUserPassword {
			return LoginAccResponse{}, errors.New("Wrong password")
		}

		token, err := CreateToken(results[0]["username"].(string), results[0]["role"].(string))

		if err != nil {
			return LoginAccResponse{}, errors.New("Token Error")
		}

		loginResponse = LoginAccResponse{
			FirstName:  results[0]["firstname"].(string),
			SecondName: results[0]["secondname"].(string),
			Username:   results[0]["username"].(string),
			Email:      results[0]["email"].(string),
			Role:       results[0]["role"].(string),
			Message:    "User Logged In",
			Token:      token,
			Success:    true,
		}

	}
	return loginResponse, nil

}

func CreateAProduct(name string, price int, description, category, defaultcurrency string) (string, error) {

	url := SupaURL + "/products"
	var finalResponse string
	productData := map[string]any{
		"name":            name,
		"price":           price,
		"description":     description,
		"defaultcurrency": defaultcurrency,
		"category":        category,
	}

	data, err := json.Marshal(productData)

	if err != nil {
		return "", errors.New("Json error")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", errors.New("Request error")
	}

	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return "", errors.New("Response error")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		return "", errors.New("Supabase Error" + string(body))
	}
	if response.StatusCode == http.StatusCreated {
		finalResponse = "Product Created in the db"
	}

	return finalResponse, nil
}

func GetAllProducts() ([]GetProduct, error) {

	url := fmt.Sprintf("%s/products?select=*", SupaURL)
	req, err := http.NewRequest("GET", url, nil)
	var result []GetProduct

	if err != nil {
		return []GetProduct{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []GetProduct{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []GetProduct
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []GetProduct{}, errors.New("Product not found in the db")

	} else {
		result = results
	}

	return result, nil
}

func GetSingleProduct(id string) ([]GetProduct, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s&select=*", SupaURL, id)
	req, err := http.NewRequest("GET", url, nil)
	var result []GetProduct

	if err != nil {
		return []GetProduct{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []GetProduct{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []GetProduct
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []GetProduct{}, errors.New("No products found")

	} else {
		result = results
	}

	return result, nil
}

func UpdateSingleProduct(id string, newName string, newPrice int) ([]GetProduct, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s", SupaURL, id)

	var result []GetProduct

	payload := map[string]interface{}{
		"name":  newName,
		"price": newPrice,
	}

	data, err := json.Marshal(payload)

	if err != nil {
		return []GetProduct{}, errors.New("Json error")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))

	if err != nil {
		return []GetProduct{}, errors.New("Request error")
	}

	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []GetProduct{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []GetProduct
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []GetProduct{}, errors.New("No products found")

	} else {
		result = results
	}

	return result, nil
}
