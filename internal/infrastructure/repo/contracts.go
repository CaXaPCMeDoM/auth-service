package repo

import (
	"auth-service/internal/entity"
	"context"
	"github.com/google/uuid"
)

type (
	UserRepo interface {
		GetEmailByID(ctx context.Context, userID uuid.UUID) (string, error)
	}
	SessionRepo interface {
		CreateSession(ctx context.Context, session entity.Session) error
		GetByJwtID(ctx context.Context, jwtID uuid.UUID) (*entity.Session, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}
	Transactor interface {
		WithinTransaction(context.Context, func(ctx context.Context) error) error
	}
)
