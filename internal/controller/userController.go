package controller

import (
	"context"
	"ivanjabrony/refstudy/internal/model/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserCotroller struct {
	userService UserUsecase
	validator   *validator.Validate
}

type UserUsecase interface {
	CreateUser(ctx context.Context, dto *dto.CreateUserDto) (*dto.UserDto, error)

	GetUserById(ctx context.Context, id int32) (*dto.UserDto, error)

	GetAllUsers(ctx context.Context) ([]dto.UserDto, error)

	UpdateUser(context.Context, *dto.UpdateUserDto) error

	DeleteUserById(context.Context, int32) error
}

func NewUserController(userService UserUsecase, validator *validator.Validate) *UserCotroller {
	return &UserCotroller{
		userService: userService,
		validator:   validator}
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  returning user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of user"
// @Success      200 {object} dto.UserDto
// @Failure      400 {object} dto.BadResponseDto
// @Router       /users/{id} [get]
func (pc *UserCotroller) GetUser(c *gin.Context) {
	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No iD provided"})
		return
	}

	parsedId64, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ID"})
		return
	}

	parsedId := int32(parsedId64)
	user, err := pc.userService.GetUserById(c.Request.Context(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUser godoc
// @Summary      Get all users with pagination
// @Description  returning users with pagination
// @Tags         user
// @Accept       json
// @Produce      json
// @Param page query int false "Page number (starting from 1)" default(1)
// @Param page_size query int false "Amount of items on the page" default(10) minimum(1) maximum(100)
// @Success      200 {object} dto.PaginatedUsersDto
// @Failure      400 {object} dto.BadResponseDto
// @Router       /users [get]
func (pc *UserCotroller) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, err := pc.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info"})
		return
	}

	total := len(users)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	response := dto.PaginatedUsersDto{
		Data:       users[offset:min(len(users), offset+pageSize)],
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	// err = pc.validator.StructCtx(c.Request.Context(), response)
	// if err != nil {
	// 	for _, e := range err.(validator.ValidationErrors) {
	// 		fmt.Println(
	// 			"Error in field:", e.Field(),
	// 			"Rule abrupted:", e.Tag(),
	// 			"Current value:", e.Value(),
	// 		)
	// 	}
	// }
	c.JSON(http.StatusOK, response)
}

// CreateUser godoc
// @Summary     Create user
// @Description Creates new user
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateUserDto true "User data"
// @Success     204 "Creating Success"
// @Failure     400 {object} dto.BadResponseDto
// @Router      /users [post]
func (pc *UserCotroller) CreateUser(c *gin.Context) {
	var createDto dto.CreateUserDto

	err := c.BindJSON(&createDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data for creating"})
		return
	}

	id, err := pc.userService.CreateUser(c.Request.Context(), &createDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, id)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Updates existing user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body dto.UpdateUserDto true "Updated data"
// @Success      204 "Update success"
// @Failure      400 {object} dto.BadResponseDto
// @Router       /users [put]
func (pc *UserCotroller) UpdateUser(c *gin.Context) {
	var updateDto dto.UpdateUserDto

	err := c.BindJSON(&updateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data for updating"})
		return
	}

	err = pc.userService.UpdateUser(c.Request.Context(), &updateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	c.JSON(http.StatusOK, updateDto.Id)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Deletes user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      204 "Delete success"
// @Failure      400 {object} dto.BadResponseDto
// @Router       /users/{id} [delete]
func (pc *UserCotroller) DeleteUserById(c *gin.Context) {
	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No iD provided"})
		return
	}

	parsedId64, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ID"})
		return
	}

	parsedId := int32(parsedId64)
	err = pc.userService.DeleteUserById(c.Request.Context(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user data"})
		return
	}

	c.JSON(http.StatusOK, id)
}
