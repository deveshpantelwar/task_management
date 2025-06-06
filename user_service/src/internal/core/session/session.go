package session

type RegisterResponse struct {
	UID      int64  `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
