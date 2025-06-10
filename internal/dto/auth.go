package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  LoginUser `json:"user"`
}

type LoginUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
