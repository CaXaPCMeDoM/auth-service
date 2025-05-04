package jwt

import (
	"auth-service/pkg/auth"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net"
	"time"
)

var (
	ErrEmptySecret = errors.New("secret cannot be empty")
	ErrNegativeTTL = errors.New("ttl must be positive")
)

type Manager struct {
	secretKey  []byte
	ttlAccess  time.Duration
	ttlRefresh time.Duration
}

func New(
	secretKey []byte,
	ttlAccess time.Duration,
	ttlRefresh time.Duration,
) (*Manager, error) {
	if len(secretKey) == 0 {
		return nil, ErrEmptySecret
	}

	if ttlAccess <= 0 || ttlRefresh <= 0 {
		return nil, ErrNegativeTTL
	}

	return &Manager{
		secretKey:  secretKey,
		ttlAccess:  ttlAccess,
		ttlRefresh: ttlRefresh,
	}, nil
}

func (m *Manager) GenerateAccess(userID uuid.UUID, clientIP net.IP) (string, error) {
	now := time.Now()
	exp := now.Add(m.ttlAccess)

	claim := &auth.Claim{
		ClientIP: clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)

	return token.SignedString(m.secretKey)
}

func (m *Manager) Parse(tokenStr string) (*auth.Claim, error) {
	trueAlgName := []string{jwt.SigningMethodHS512.Alg()}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&auth.Claim{},
		func(token *jwt.Token) (interface{}, error) {
			return m.secretKey, nil
		},
		jwt.WithValidMethods(trueAlgName),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*auth.Claim)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (m *Manager) GenerateRefresh() (refreshToken string, expired time.Time) {
	return uuid.NewString(), time.Now().Add(m.ttlRefresh)
}
