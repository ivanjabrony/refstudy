package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ivanjabrony/refstudy/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{db}
}

func (repo UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Insert("users").
		Columns("username", "email", "password").
		Values(
			user.Username,
			user.Email,
			user.Password,
		).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&user.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not inserted: %w", err)
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return user, nil
}

func (repo UserRepository) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Select("id, username", "email", "password").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	err = tx.SelectContext(ctx, &user, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not inserted: %w", err)
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}

func (repo UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Select("id, username", "email", "password").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var users []model.User

	err = tx.SelectContext(ctx, &users, query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return users, nil
}

func (repo UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Where(squirrel.Eq{"id": user.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.Id)
	}
	return nil
}
func (repo UserRepository) DeleteUserById(ctx context.Context, id int32) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
