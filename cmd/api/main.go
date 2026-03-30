package main

import (
	"net/http"

	"github.com/idorocodes/ecom/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1", internal.HomePath)
	mux.HandleFunc("GET /api/v1/health", internal.HealthStatus)
	mux.HandleFunc("POST /api/v1/createuser", internal.CreateAccount)
	mux.HandleFunc("POST /api/v1/loginuser", internal.LoginAccount)
	mux.Handle("GET /api/v1/dashboard", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.DashboardPath)))
	mux.Handle("POST /api/v1/createproduct", internal.AuthMiddleware("admin")(http.HandlerFunc(internal.CreateProduct)))
	mux.Handle("GET /api/v1/getproducts", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.GetProducts)))
	mux.Handle("GET /api/v1/getproduct", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.GetOneProduct)))
	mux.Handle("PATCH /api/v1/updateproduct", internal.AuthMiddleware("admin")(http.HandlerFunc(internal.UpdateAProduct)))
	mux.Handle("DELETE /api/v1/deleteproduct", internal.AuthMiddleware("admin")(http.HandlerFunc(internal.DeleteProduct)))
	mux.Handle("POST /api/v1/createorder", internal.AuthMiddleware("admin", "customer")(http.HandlerFunc(internal.CreateOrder)))
	mux.Handle("GET /api/v1/getmetrics", internal.AuthMiddleware("admin")(http.HandlerFunc(internal.Metrics)))
	fmt.Println("Server started on http://localhost:8080!")

	http.ListenAndServe(":8080", mux)

}
