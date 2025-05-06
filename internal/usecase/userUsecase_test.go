package usecase_test

import (
	"context"
	"errors"
	"ivanjabrony/refstudy/internal/logger"
	"ivanjabrony/refstudy/internal/model"
	"ivanjabrony/refstudy/internal/model/dto"
	"ivanjabrony/refstudy/internal/usecase"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserStorage) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserStorage) GetAllUsers(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *mockUserStorage) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(1)
}

func (m *mockUserStorage) DeleteUserById(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(1)
}

func TestNewUserUsecase(t *testing.T) {
	testcases := []struct {
		name       string
		repository usecase.UserRepository
		err        error
		isNil      bool
	}{
		{
			name:       "success",
			repository: new(mockUserStorage),
			err:        nil,
			isNil:      false,
		},
		{
			name:       "nil user storage",
			repository: nil,
			err:        errors.New("nil values in UserUsecase constructor"),
			isNil:      true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			service, err := usecase.NewUserUsecase(testcase.repository, &logger.MyLogger{})
			if err != nil {
				require.Error(t, testcase.err)
				require.ErrorContains(t, err, testcase.err.Error())
			}
			if testcase.isNil {
				require.Nil(t, service)
			} else {
				require.NotNil(t, service)
			}
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctx := context.Background()
	password := "password"
	user_id := int32(1)
	storagedUser := &model.User{
		Id:       1,
		Username: "ivan",
		Email:    "test@example.com",
		Password: password,
	}
	expectedUser := &dto.UserDto{
		Id:       1,
		Username: "ivan",
		Email:    "test@example.com",
		Password: password,
	}

	storage1 := new(mockUserStorage)
	storage2 := new(mockUserStorage)
	storage1.On("GetUserById", ctx, mock.Anything).Return(storagedUser, nil)
	storage2.On("GetUserById", ctx, mock.Anything).Return(&model.User{}, errors.New(""))

	for _, testcase := range []struct {
		name         string
		storage      *mockUserStorage
		expectedUser *dto.UserDto
		err          error
	}{
		{
			name:         "success",
			storage:      storage1,
			expectedUser: expectedUser,
			err:          nil,
		},
		{
			name:         "user exists",
			storage:      storage2,
			expectedUser: nil,
			err:          errors.New(""),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// arrange
			service, _ := usecase.NewUserUsecase(testcase.storage, &logger.MyLogger{})

			// act
			user, err := service.GetUserById(ctx, user_id)

			// assert
			if err != nil {
				require.Error(t, testcase.err)
				require.Error(t, err)
			}
			require.Equal(t, user, testcase.expectedUser)
			testcase.storage.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_CreateUser(t *testing.T) {
	ctx := context.Background()
	password := "password"
	payload := dto.CreateUserDto{
		Username: "ivan",
		Email:    "test@example.com",
		Password: password,
	}
	userToStorage := &model.User{
		Username: "ivan",
		Email:    "test@example.com",
		Password: password,
	}
	storagedUser := &model.User{
		Id:       1,
		Username: "ivan",
		Email:    "test@example.com",
		Password: password,
	}

	for _, testcase := range []struct {
		name         string
		storageSetup func(*mockUserStorage)
		expectedUser *dto.UserDto
		err          error
	}{
		{
			name: "success",
			storageSetup: func(m *mockUserStorage) {
				m.On("CreateUser", ctx, userToStorage).Return(storagedUser, nil)
			},
			expectedUser: &dto.UserDto{
				Id:       1,
				Username: "ivan",
				Email:    "test@example.com",
				Password: password,
			},
			err: nil,
		},
		{
			name: "creation error",
			storageSetup: func(m *mockUserStorage) {
				m.On("CreateUser", ctx, userToStorage).Return(&model.User{}, errors.New("error"))
			},
			expectedUser: nil,
			err:          errors.New(""),
		}} {
		t.Run(testcase.name, func(t *testing.T) {
			// arrange
			storage := new(mockUserStorage)
			testcase.storageSetup(storage)
			service, _ := usecase.NewUserUsecase(
				storage,
				&logger.MyLogger{},
			)

			// act
			user, err := service.CreateUser(ctx, &payload)

			// assert
			if err != nil {
				require.Error(t, testcase.err)
				require.Error(t, err)
			}
			require.Equal(t, user, testcase.expectedUser)
			storage.AssertExpectations(t)
		})
	}
}
