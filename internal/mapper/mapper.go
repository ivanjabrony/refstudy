package mapper

import (
	"ivanjabrony/refstudy/internal/model"
	"ivanjabrony/refstudy/internal/model/dto"
)

func MapFromCreateUserDto(dto *dto.CreateUserDto) *model.User {
	if dto != nil {
		return &model.User{
			Username: dto.Username,
			Email:    dto.Email,
			Password: dto.Password,
		}
	}

	return nil
}

func MapFromUserDto(dto *dto.UserDto) *model.User {
	if dto != nil {
		return &model.User{
			Id:       dto.Id,
			Username: dto.Username,
			Email:    dto.Email,
			Password: dto.Password,
		}
	}

	return nil
}

func MapFromUpdateUserDto(dto *dto.UpdateUserDto) *model.User {
	var updUsername, updEmail, updPassword string
	if dto.Username != nil {
		updUsername = *dto.Username
	}
	if dto.Email != nil {
		updEmail = *dto.Email
	}
	if dto.Password != nil {
		updPassword = *dto.Password
	}
	if dto != nil {
		return &model.User{
			Id:       dto.Id,
			Username: updUsername,
			Email:    updEmail,
			Password: updPassword,
		}
	}

	return nil
}

func MapToUserDto(model *model.User) *dto.UserDto {
	if model != nil {
		return &dto.UserDto{
			Id:       model.Id,
			Username: model.Username,
			Email:    model.Email,
			Password: model.Password,
		}
	}

	return nil
}

func MapToManyUserDto(models ...model.User) []dto.UserDto {
	dtos := make([]dto.UserDto, len(models))
	for i, v := range models {
		dtos[i] = *MapToUserDto(&v)
	}

	return dtos
}

func MapFromManyUserDto(dtos ...dto.UserDto) []model.User {
	models := make([]model.User, len(dtos))
	for i, v := range dtos {
		models[i] = *MapFromUserDto(&v)
	}

	return models
}
