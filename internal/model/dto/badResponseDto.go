package dto

type BadResponseDto struct {
	Response string `json:"error" example:"Server error"`
}
