package persistance

import (
	"context"
	"fmt"
	errors "task-management/user-service/src/internal/core/errors"
	"task-management/user-service/src/internal/core/session"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(db *Database) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) IsEmailOsUserNameTaken(ctx context.Context, email, username string) (bool, error) {
	var count int
	err := ur.db.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2`, email, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can not check user : %v", err)
	}
	return count > 0, nil
}

func (ur *UserRepo) RegisterUser(ctx context.Context, u *session.RegisterResponse) error {

	var userID int64

	err := ur.db.db.QueryRow(`
	INSERT INTO users
	(username, email, password) VALUES ($1, $2, $3) RETURNING uid`,
		u.UserName, u.Email, u.Password,
	).Scan(&userID)

	if err != nil {
		return fmt.Errorf("%v : failed inserting user", errors.ErrInternalServer)
	}

	return nil
}

func (ur *UserRepo) GetUserByUsername(ctx context.Context, username string) (*session.RegisterResponse, error) {
	var user session.RegisterResponse
	err := ur.db.db.QueryRow(`SELECT uid, username, email, password FROM users WHERE username = $1`, username).
		Scan(&user.UID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, uid int) (*session.RegisterResponse, error) {
	var user session.RegisterResponse

	query := `SELECT uid, username, email, password_hash FROM users WHERE uid = $1`
	err := ur.db.db.QueryRowContext(ctx, query, uid).Scan(&user.UID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user *session.RegisterResponse) error {

	fmt.Println("request details : ", user)
	fmt.Printf(">>> Incoming update request: username=%s, email=%s\n", user.UserName, user.Email)

	result, err := ur.db.db.ExecContext(ctx,
		`UPDATE users SET email = $1 WHERE username = $2`,
		user.Email, user.UserName,
	)

	var newData session.RegisterResponse

	err = ur.db.db.QueryRowContext(ctx, `SELECT * FROM users WHERE username=$1`, user.UserName).Scan(&newData.UID, &newData.UserName, &newData.Email, &newData.Password)

	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)

		return err
	}

	// fmt.Printf("Trying to update user: username=%s, email=%s, password=%s\n", user.UserName, user.Email, user.Password)
	fmt.Println("New user details : ", newData)

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
