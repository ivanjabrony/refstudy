package repository_test

import (
	"context"
	"ivanjabrony/refstudy/internal/logger"
	"ivanjabrony/refstudy/internal/model"
	"ivanjabrony/refstudy/internal/repository"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPgxPool struct {
	mock.Mock
	pgxpool.Pool
}

func TestNewUserStorage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pool := &pgxpool.Pool{}
		storage, err := repository.NewUserRepository(pool, &logger.MyLogger{})
		require.NoError(t, err)
		require.NotNil(t, storage)
	})

	t.Run("nil pool", func(t *testing.T) {
		storage, err := repository.NewUserRepository(nil, &logger.MyLogger{})
		require.Error(t, err)
		require.Nil(t, storage)
	})
}

func TestShouldGetUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	repo, err := repository.NewUserRepository(mock, &logger.MyLogger{})
	require.NoError(t, err)

	var id int32 = 1
	rs := pgxmock.
		NewRows([]string{"id", "username", "email", "password"}).
		AddRow(id, "ivan", "123@example.com", "12345678")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, username, email, password FROM users WHERE id = \\$1").WithArgs(id).WillReturnRows(rs)
	mock.ExpectCommit()

	user, err := repo.GetUserById(context.Background(), id)
	require.NoError(t, err)

	require.Equal(t, user.Id, id)
	require.Equal(t, user.Username, "ivan")
	require.Equal(t, user.Email, "123@example.com")
	require.Equal(t, user.Password, "12345678")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldCreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	repo, err := repository.NewUserRepository(mock, &logger.MyLogger{})
	require.NoError(t, err)

	var id int32 = 1
	rs := pgxmock.
		NewRows([]string{"id"}).
		AddRow(id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WithArgs("ivan", "123@example.com", "12345678").WillReturnRows(rs)
	mock.ExpectCommit()

	user, err := repo.CreateUser(context.Background(),
		&model.User{
			Id:       0,
			Username: "ivan",
			Email:    "123@example.com",
			Password: "12345678"})
	require.NoError(t, err)

	require.Equal(t, user.Id, id)
	require.Equal(t, user.Username, "ivan")
	require.Equal(t, user.Email, "123@example.com")
	require.Equal(t, user.Password, "12345678")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
