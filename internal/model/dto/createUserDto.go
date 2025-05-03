package dto

type CreateUserDto struct {
	Username string `json:"username" example:"Ivan" binding:"required"`
	Email    string `json:"email" example:"123@example.com" binding:"required"`
	Password string `json:"password" example:"12345678" binding:"required"`
}
