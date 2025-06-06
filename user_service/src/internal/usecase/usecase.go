package usecase

import (
	"context"
	"task-management/user-service/src/internal/core/errors"
	"task-management/user-service/src/internal/core/session"
	"task-management/user-service/src/internal/core/user"
	"task-management/user-service/src/pkg"
)

type UserService struct {
	userRepo  user.Repository
	secretKey string
}

func NewUserService(userRepo user.Repository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		secretKey: jwtSecret,
	}
}

type RegisterInput struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserService) RegisterUser(ctx context.Context, input RegisterInput) error {

	isTaken, err := u.userRepo.IsEmailOsUserNameTaken(ctx, input.Email, input.UserName)
	if err != nil {
		return err
	}
	if isTaken {
		return errors.ErrEmailOrUsernameTaken
	}

	hashedPassword, err := pkg.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := &session.RegisterResponse{
		UserName: input.UserName,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	return u.userRepo.RegisterUser(ctx, user)
}

//---------------------------------------------------

type LoginInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

func (u *UserService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Fetch user by username
	userData, err := u.userRepo.GetUserByUsername(ctx, input.UserName)
	if err != nil {
		return nil, err
	}
	if userData == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check password
	if !pkg.CheckPassword(userData.Password, input.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := pkg.GenerateJWT(u.secretKey, int(userData.UID))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &LoginOutput{Token: token}, nil
}

type UserOutput struct {
	UID      int    `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserService) GetUserByUsername(ctx context.Context, username string) (*session.RegisterResponse, error) {
	return u.userRepo.GetUserByUsername(ctx, username)
}

type UpdateUserInput struct {
	UserName    string `json:"username"`
	NewEmail    string `json:"email"`
	NewPassword string `json:"password"`
}

func (u *UserService) UpdateUser(ctx context.Context, input UpdateUserInput) error {

	existingUser, err := u.userRepo.GetUserByUsername(ctx, input.UserName)
	if err != nil {
		return errors.ErrUserNotFound
	}

	if input.NewEmail != "" {
		existingUser.Email = input.NewEmail
	}
	if input.NewPassword != "" {
		hashed, err := pkg.HashPassword(input.NewPassword)
		if err != nil {
			return err
		}

		existingUser.Password = hashed
	}

	return u.userRepo.UpdateUser(ctx, existingUser)
}
