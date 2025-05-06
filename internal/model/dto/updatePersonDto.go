package dto

type UpdateUserDto struct {
	Id       int32   `json:"id" example:"1" binding:"required" validate:"required,gt=0"`
	Username *string `json:"username" example:"Ivan" validate:"printascii"`
	Email    *string `json:"email" example:"123@example.com" validate:"email"`
	Password *string `json:"password" example:"12345678"`
}
