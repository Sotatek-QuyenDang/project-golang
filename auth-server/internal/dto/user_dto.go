package dto

type LoginRequest struct {
	UserName       string `json:"user_name" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
}
type LoginResponse struct {
	Token string `json:"token"`
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
