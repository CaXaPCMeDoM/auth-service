package persistence

import (
	"auth-service/internal/entity"
	"auth-service/pkg/postgres"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

type SessionRepo struct {
	*postgres.Postgres
}

func NewSessionRepo(pg *postgres.Postgres) *SessionRepo {
	return &SessionRepo{pg}
}

func (r *SessionRepo) CreateSession(ctx context.Context, session entity.Session) error {
	sql, args, err := r.Builder.
		Insert("session").
		Columns(
			"id",
			"refresh_token_hash",
			"jwt_id",
			"user_ip",
			"expired_at",
			"created_at",
			"user_id",
		).
		Values(
			session.ID,
			session.RefreshTokenHash,
			session.Jti,
			session.UserIP.String(),
			session.ExpiredAt,
			session.CreatedAt,
			session.UserID,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo - CreateSession - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("SessionRepo - CreateSession - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) GetByJwtID(ctx context.Context, jwtID uuid.UUID) (*entity.Session, error) {
	sql, args, err := r.Builder.
		Select("id", "refresh_token_hash", "jwt_id", "user_ip", "expired_at", "user_id").
		From("session").
		Where(squirrel.And{
			squirrel.Eq{"jwt_id": jwtID},
			squirrel.Gt{"expired_at": time.Now().UTC()},
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SessionRepo - GetSessionByJwtID - r.Builder: %w", err)
	}

	var s entity.Session

	err = r.QueryRow(ctx, sql, args...).Scan(
		&s.ID, &s.RefreshTokenHash, &s.Jti, &s.UserIP, &s.ExpiredAt, &s.UserID,
	)

	if err != nil {
		return nil, entity.ErrSessionNotFound
	}

	return &s, nil
}

func (r *SessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := r.Builder.
		Delete("session").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("SessionRepo - Delete - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)

	if err != nil {
		return entity.ErrSessionNotFound
	}

	return nil
}
