package main

import (
	"fmt"
	"net/http"
	"github.com/idorocodes/ecom/internal"
)

func main() {

	http.HandleFunc("/", internal.GetHomePath)
	fmt.Println("Server staring!")
	http.ListenAndServe(":8080", nil)

}
  
 
  