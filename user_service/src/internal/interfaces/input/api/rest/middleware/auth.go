// package middleware

// import (
// 	"net/http"
// 	"task-management/user-service/src/internal/adaptors/persistance"
// )

// type AuthMiddleware struct {
// 	userRepo persistance.UserRepo
// }

// func NewAuthMiddleware(userRepo persistance.UserRepo) *AuthMiddleware {
// 	return &AuthMiddleware{userRepo: userRepo}
// }

// func (m *AuthMiddleware) Validate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("access_token")
// 		if err != nil || cookie.Value == "" {
// 			http.Error(w, "Unauthorized - No token cookie", http.StatusUnauthorized)
// 			return
// 		}

// 		valid, err := m.userRepo.GetUserByToken(cookie.Value)
// 		if err != nil || !valid {
// 			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: secret}
}

func (m *AuthMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token") // Must match the cookie name you set
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
