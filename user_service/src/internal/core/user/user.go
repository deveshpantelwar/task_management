package user

import (
	"context"
	"task-management/user-service/src/internal/core/session"
)

type User struct {
	UID      int    `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Repository interface {
	RegisterUser(ctx context.Context, user *session.RegisterResponse) error
	IsEmailOsUserNameTaken(ctx context.Context, email, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*session.RegisterResponse, error)
	GetUserByID(ctx context.Context, uid int) (*session.RegisterResponse, error)
	UpdateUser(ctx context.Context, user *session.RegisterResponse) error
}
