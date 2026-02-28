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
	mux.Handle("GET /dashboard", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.DashboardPath)))
	mux.Handle("POST /createproduct", internal.AuthMiddleware("admin")(http.HandlerFunc(internal.CreateProduct)))
	mux.Handle("GET /getProduct", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.GetProducts)))
	fmt.Println("Server started on http://localhost:8080!")

	http.ListenAndServe(":8080", mux)

}
