package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(allowedRole ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			isValid, userId, userRole, err := VerifyToken(token)

			if len(userId) == 0 || len(userRole) == 0 {
				http.Error(w, "Token Error", http.StatusUnauthorized)
				return
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
				return
			}
			if !isValid {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
				return
			}

			isAllowed := false

			for _, role := range allowedRole {
				if role == userRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden: Insufficient Permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userId, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
