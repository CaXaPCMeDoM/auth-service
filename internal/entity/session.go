package entity

import (
	"errors"
	"github.com/google/uuid"
	"net"
	"time"
)

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrBadParse            = errors.New("bad parse to session entity")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrIPNotFound          = errors.New("ip not found")
	ErrGenerateAccess      = errors.New("error generate access token")
	ErrGenerateHash        = errors.New("error generate hash")
	ErrTransaction         = errors.New("error with transaction")
	ErrCreateSession       = errors.New("error create session")
	ErrUnauthorized        = errors.New("unauthorized")
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	UserIP           net.IP
	Jti              uuid.UUID
	RefreshTokenHash string
	ExpiredAt        time.Time
	CreatedAt        time.Time
}
