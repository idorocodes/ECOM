package main

import (
	"fmt"
	"net/http"

	"github.com/idorocodes/ecom/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(" /", internal.HomePath)
	mux.HandleFunc(" /createuser", internal.CreateAccount)
	mux.HandleFunc("/loginuser", internal.LoginAccount)
	mux.Handle("/dashboard",  internal.AuthMiddleware("customer")(http.HandlerFunc(internal.DashboardPath)))
	fmt.Println("Server started on http://localhost:8080!")

	http.ListenAndServe(":8080", mux)

}
