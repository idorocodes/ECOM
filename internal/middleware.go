package internal

import (
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		isValid, token, err := VerifyToken(token)

		if err != nil {
			http.Error(w, fmt.Sprintf("%v",err), http.StatusUnauthorized)
			return
		}
		if !isValid {
			http.Error(w,  fmt.Sprintf("%v",err), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
