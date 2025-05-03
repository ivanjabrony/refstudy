package dto

type UserDto struct {
	Id       int32  `json:"id" example:"1"`
	Username string `json:"username" example:"Ivan"`
	Email    string `json:"email" example:"123@example.com"`
	Password string `json:"password" example:"12345678"`
}
