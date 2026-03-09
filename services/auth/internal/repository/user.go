package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iamsainadh/neovault/auth/internal/app_error"
	"github.com/iamsainadh/neovault/auth/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (*model.CreatedUser, error)
	LoginUser(ctx context.Context, input *model.LoginUserInput) (*model.LoginUserResponse, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{db: pool}
}

func (ur *userRepository) CreateUser(ctx context.Context, u model.User) (*model.CreatedUser, error) {

	currCxt, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	var query = `INSERT INTO USERS(name,email,password) VALUES ($1,$2,$3) returning id,name,email`

	err := ur.db.QueryRow(currCxt, query, u.Name, u.Email, u.Password).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
	)

	if err != nil {
		return &model.CreatedUser{}, mapDBError(err, "")
	}

	return &model.CreatedUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, err
}

func mapDBError(err error, op string) error {

	if errors.Is(err, errNoRows()) {
		return fmt.Errorf("%s: %w", op, app_error.ErrUserNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			if pgErr.ConstraintName == "users_email_key" {
				return fmt.Errorf("%s: %w", op, app_error.ErrEmailAlreadyExist)
			}
		}
	}

	return fmt.Errorf("%s: %w", op, err)
}

func errNoRows() error {

	return pgxNoRowsErr{}
}

type pgxNoRowsErr struct{}

func (pgxNoRowsErr) Error() string { return "no rows in result set" }

// User Login
func (r *userRepository) LoginUser(ctx context.Context, input *model.LoginUserInput) (*model.LoginUserResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	var user model.LoginUserResponse

	query := `SELECT id,name,email,password FROM USERS where email = $1 and is_active=true and deleted_at is null LIMIT 1;`

	err := r.db.QueryRow(currCtx, query, input.Email).Scan(
		&user.UserId,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return &model.LoginUserResponse{}, mapDBError(err, "")
	}

	return &model.LoginUserResponse{
		UserId:   user.UserId,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
