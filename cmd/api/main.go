package main

import (
	"fmt"
	"net/http"
	"github.com/idorocodes/ecom/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", internal.HomePath)
	mux.HandleFunc("POST /users", internal.CreateAccount)

	fmt.Println("Server started on http://localhost:8080!")

	http.ListenAndServe(":8080", mux)

}
