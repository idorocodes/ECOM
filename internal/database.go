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
	SupaURL = ""
	SupaKey = ""
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

func GetAllUsers() ([]AllUsers, error) {

	url := fmt.Sprintf("%s/users?select=*", SupaURL)
	req, err := http.NewRequest("GET", url, nil)
	var result []AllUsers

	if err != nil {
		return []AllUsers{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []AllUsers{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []AllUsers
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []AllUsers{}, errors.New("No users found in the db")

	} else {
		result = results
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

func CreateAProduct(name string, price int, description, category, defaultcurrency, status string) (string, error) {

	url := SupaURL + "/products"
	var finalResponse string
	productData := map[string]any{
		"name":            name,
		"price":           price,
		"description":     description,
		"status":          status,
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

func GetAllProducts() ([]Product, error) {

	url := fmt.Sprintf("%s/products?select=*", SupaURL)
	req, err := http.NewRequest("GET", url, nil)
	var result []Product

	if err != nil {
		return []Product{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []Product{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []Product
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []Product{}, errors.New("Product not found in the db")

	} else {
		result = results
	}

	return result, nil
}

func GetSingleProduct(id string) (GetProduct, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s&select=*", SupaURL, id)
	req, err := http.NewRequest("GET", url, nil)
	var result GetProduct

	if err != nil {
		return GetProduct{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return GetProduct{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []GetProduct
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return GetProduct{}, errors.New("No products found")

	} else {
		result = results[0]
	}

	return result, nil
}
func UpdateSingleProduct(id string, newName string, newPrice int) ([]GetProduct, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s", SupaURL, id)

	payload := map[string]interface{}{
		"name":  newName,
		"price": newPrice,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+SupaKey)
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase error: status %d", resp.StatusCode)
	}

	var results []GetProduct
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(results) == 0 {
		return nil, errors.New("no products matched the ID")
	}

	return results, nil
}
func DeleteSingleProduct(id string) (bool, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s", SupaURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var deletedItems []GetProduct
	json.NewDecoder(resp.Body).Decode(&deletedItems)

	if len(deletedItems) == 0 {
		return false, errors.New("nothing was deleted (ID not found)")
	}

	return true, nil
}

func CreateAnOrder(id, user_id, address string) (string, int, error) {
	url := fmt.Sprintf("%s/products?id=eq.%s", SupaURL, id)

	payload := map[string]interface{}{
		"status": "ordered",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return "", 0, fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+SupaKey)
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", 0, fmt.Errorf("supabase error: status %d", resp.StatusCode)
	}

	var results []GetProduct
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", 0, fmt.Errorf("decode error: %w", err)
	}

	if len(results) == 0 {
		return "", 0, errors.New("no products matched the ID")
	}

	orderurl := SupaURL + "/orders"

	var finalResponse string
	productData := map[string]any{
		"product_id": results[0].Id,
		"user_id":    user_id,
		"address":    address,
	}

	orderdata, err := json.Marshal(productData)

	if err != nil {
		return "", 0, errors.New("Json error")
	}

	orderreq, err := http.NewRequest("POST", orderurl, bytes.NewBuffer(orderdata))
	if err != nil {
		return "", 0, errors.New("Request error")
	}

	orderreq.Header.Set("apiKey", SupaKey)
	orderreq.Header.Set("Authorization", "Bearer "+SupaKey)
	orderreq.Header.Set("Content-Type", "application/json")
	orderreq.Header.Set("Prefer", "return=minimal")

	orderclient := &http.Client{}

	orderresponse, err := orderclient.Do(orderreq)

	if err != nil {
		return "", 0, errors.New("Response error")
	}

	defer orderresponse.Body.Close()

	if orderresponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(orderresponse.Body)
		return "", 0, errors.New("Supabase Error" + string(body))
	}
	if orderresponse.StatusCode == http.StatusCreated {
		finalResponse = "Order Created in the db"
	}

	return finalResponse, orderresponse.StatusCode, nil

}

func GetAllOrders() ([]AllOrders, error) {

	url := fmt.Sprintf("%s/orders?select=*", SupaURL)
	req, err := http.NewRequest("GET", url, nil)
	var result []AllOrders

	if err != nil {
		return []AllOrders{}, errors.New("Request error")
	}
	req.Header.Set("apiKey", SupaKey)
	req.Header.Set("Authorization", "Bearer "+SupaKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return []AllOrders{}, errors.New("Response error")
	}

	defer response.Body.Close()

	var results []AllOrders
	json.NewDecoder(response.Body).Decode(&results)

	if len(results) == 0 {
		return []AllOrders{}, errors.New("No orders found in the db")

	} else {
		result = results
	}

	return result, nil
}
