package middleware

import (
	"context"
	"net/http"
	"strings"
	"task-management/user-service/src/pkg"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware validates JWT and sets userID in context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			pkg.Error(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token := tokenParts[1]

		// Validate token and get claims
		claims, err := pkg.VerifyJWT(token)
		if err != nil {
			pkg.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Extract userID from claims and put into context
		userID := claims.UID

		ctx := context.WithValue(r.Context(), UserIDKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
