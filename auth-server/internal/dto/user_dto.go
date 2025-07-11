package dto

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
type CreateUserRequest struct {
	UserName       string `json:"user_name" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
	Role           string `json:"role" binding:"required"`
}
type UpdateUserRequest struct {
	HashedPassword string `json:"hashed_password"`
	Role           string `json:"role"`
}
type UserResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
}
