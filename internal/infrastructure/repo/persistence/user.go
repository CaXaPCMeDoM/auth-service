package persistence

import (
	"auth-service/pkg/postgres"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetEmailByID(ctx context.Context, userID uuid.UUID) (string, error) {
	sql, args, err := r.Builder.
		Select("email").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("UserRepo - GetEmailByID - r.Builder: %w", err)
	}

	var email string

	err = r.QueryRow(ctx, sql, args...).Scan(&email)

	if err != nil {
		return "", fmt.Errorf("UserRepo - GetEmailByID - rosw.Scan: %w", err)
	}

	return email, nil
}
