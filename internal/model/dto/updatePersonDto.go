package dto

type UpdateUserDto struct {
	Id       int32   `json:"id" example:"1" binding:"required" `
	Username *string `json:"username" example:"Ivan"`
	Email    *string `json:"email" example:"123@example.com"`
	Password *string `json:"password" example:"12345678"`
}
