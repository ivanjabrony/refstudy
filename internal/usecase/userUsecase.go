package usecase

import (
	"context"
	"ivanjabrony/refstudy/internal/mapper"
	"ivanjabrony/refstudy/internal/model"
	"ivanjabrony/refstudy/internal/model/dto"
)

type userRepository interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	GetUserById(context.Context, int32) (*model.User, error)
	GetAllUsers(context.Context) ([]model.User, error)
	UpdateUser(context.Context, *model.User) error
	DeleteUserById(context.Context, int32) error
}

type UserUsecase struct {
	userRepository
}

func NewUserUsecase(repo userRepository) UserUsecase {
	return UserUsecase{repo}
}

func (uc UserUsecase) CreateUser(ctx context.Context, dto *dto.CreateUserDto) (*dto.UserDto, error) {
	user := mapper.MapFromCreateUserDto(dto)
	user, err := uc.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return mapper.MapToUserDto(user), nil
}

func (uc UserUsecase) GetUserById(ctx context.Context, id int32) (*dto.UserDto, error) {
	user, err := uc.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.MapToUserDto(user), nil
}

func (uc UserUsecase) GetAllUsers(ctx context.Context) ([]dto.UserDto, error) {
	users, err := uc.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.MapToManyUserDto(users...), err
}

func (uc UserUsecase) UpdateUser(ctx context.Context, dto *dto.UpdateUserDto) error {
	user := mapper.MapFromUpdateUserDto(dto)
	return uc.userRepository.UpdateUser(ctx, user)
}

func (uc UserUsecase) DeleteUserById(ctx context.Context, id int32) error {
	return uc.userRepository.DeleteUserById(ctx, id)
}
