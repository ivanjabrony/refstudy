package dto

type PaginatedUsersDto struct {
	Data       []UserDto `json:"data"`
	Total      int       `json:"total" validate:"required"`
	Page       int       `json:"page" validate:"required"`
	PageSize   int       `json:"page_size" validate:"required"`
	TotalPages int       `json:"total_pages" validate:"required"`
}
