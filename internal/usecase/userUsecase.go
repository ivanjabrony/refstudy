package usecase

import (
	"context"
	"errors"
	"ivanjabrony/refstudy/internal/logger"
	"ivanjabrony/refstudy/internal/mapper"
	"ivanjabrony/refstudy/internal/model"
	"ivanjabrony/refstudy/internal/model/dto"
)

type UserRepository interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	GetUserById(context.Context, int32) (*model.User, error)
	GetAllUsers(context.Context) ([]model.User, error)
	UpdateUser(context.Context, *model.User) error
	DeleteUserById(context.Context, int32) error
}

type UserUsecase struct {
	UserRepository
	logger *logger.MyLogger
}

func NewUserUsecase(repo UserRepository, logger *logger.MyLogger) (*UserUsecase, error) {
	if repo == nil {
		return nil, errors.New("nil values in UserUsecase constructor")
	}
	return &UserUsecase{repo, logger}, nil
}

func (uc UserUsecase) CreateUser(ctx context.Context, dto *dto.CreateUserDto) (*dto.UserDto, error) {
	user := mapper.MapFromCreateUserDto(dto)
	user, err := uc.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return mapper.MapToUserDto(user), nil
}

func (uc UserUsecase) GetUserById(ctx context.Context, id int32) (*dto.UserDto, error) {
	user, err := uc.UserRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.MapToUserDto(user), nil
}

func (uc UserUsecase) GetAllUsers(ctx context.Context) ([]dto.UserDto, error) {
	users, err := uc.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.MapToManyUserDto(users...), err
}

func (uc UserUsecase) UpdateUser(ctx context.Context, dto *dto.UpdateUserDto) error {
	user := mapper.MapFromUpdateUserDto(dto)
	return uc.UserRepository.UpdateUser(ctx, user)
}

func (uc UserUsecase) DeleteUserById(ctx context.Context, id int32) error {
	return uc.UserRepository.DeleteUserById(ctx, id)
}
