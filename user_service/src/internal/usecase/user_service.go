package usecase

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// type UserService struct {
// 	repo      Repository // whatever your user repo interface is
// 	jwtSecret string
// }

// func NewUserService(r Repository, secret string) *UserService {
// 	return &UserService{
// 		repo:      r,
// 		jwtSecret: secret,
// 	}
// }

// func (u *UserService) ValidateToken(tokenStr string) (string, error) {
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(u.jwtSecret), nil
// 	})
// 	if err != nil || !token.Valid {
// 		return "", errors.New("invalid token")
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return "", errors.New("invalid token claims")
// 	}

// 	userID, ok := claims["user_id"].(string)
// 	if !ok {
// 		return "", errors.New("user_id not found in token")
// 	}

// 	return userID, nil
// }

func (u *UserService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// user_id might be float64 if stored as number
	switch id := claims["user_id"].(type) {
	case string:
		return id, nil
	case float64:
		// convert number to string
		return fmt.Sprintf("%.0f", id), nil
	default:
		return "", errors.New("user_id not found in token")
	}
}
