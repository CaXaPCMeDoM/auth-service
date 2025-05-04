package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net"
	"time"
)

type Claim struct {
	ClientIP net.IP `json:"ip"`
	jwt.RegisteredClaims
}

type (
	TokenManager interface {
		GenerateAccess(userID uuid.UUID, clientIP net.IP) (string, error)
		Parse(tokenStr string) (*Claim, error)
		GenerateRefresh() (refreshToken string, expired time.Time)
	}
)
