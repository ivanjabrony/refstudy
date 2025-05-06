package repository

import (
	"context"
	"errors"
	"fmt"
	"ivanjabrony/refstudy/internal/logger"
	"ivanjabrony/refstudy/internal/model"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type UserRepository struct {
	pool    PgxIface
	builder squirrel.StatementBuilderType
	logger  *logger.MyLogger
}

func NewUserRepository(pool PgxIface, logger *logger.MyLogger) (*UserRepository, error) {
	if pool == nil {
		return nil, errors.New("nil values in UserRepository constructor")
	}

	return &UserRepository{
		pool:    pool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		logger:  logger,
	}, nil
}

func (repo UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if rbErr := tx.Rollback(rollbackCtx); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		}
	}()

	query, args, err := repo.builder.
		Insert("users").
		Columns("username", "email", "password").
		Values(user.Username, user.Email, user.Password).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return user, nil
}

func (repo UserRepository) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if rbErr := tx.Rollback(rollbackCtx); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		}
	}()

	query, args, err := repo.builder.
		Select("id", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	err = tx.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return &user, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if rbErr := tx.Rollback(rollbackCtx); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %w", commitErr)
			}
		}
	}()

	query, args, err := repo.builder.
		Select("id", "username", "email", "password").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if rbErr := tx.Rollback(rollbackCtx); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %w", commitErr)
			}
		}
	}()

	query, args, err := repo.builder.
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Where(squirrel.Eq{"id": user.Id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.Id)
	}

	return nil
}

func (repo *UserRepository) DeleteUserById(ctx context.Context, id int32) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if rbErr := tx.Rollback(rollbackCtx); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %w", commitErr)
			}
		}
	}()

	query, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
