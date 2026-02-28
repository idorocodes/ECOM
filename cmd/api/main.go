package main

import (
	"fmt"
	"net/http"

	"github.com/idorocodes/ecom/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", internal.HomePath)
	mux.HandleFunc("GET /health", internal.HealthStatus)
	mux.HandleFunc("POST /createuser", internal.CreateAccount)
	mux.HandleFunc("POST /loginuser", internal.LoginAccount)
	mux.Handle("GET /dashboard", internal.AuthMiddleware("customer")(http.HandlerFunc(internal.DashboardPath)))
	fmt.Println("Server started on http://localhost:8080!")

	http.ListenAndServe(":8080", mux)

}
