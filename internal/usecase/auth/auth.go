package auth

import (
	"auth-service/internal/controller/http/dto"
	"auth-service/internal/entity"
	"auth-service/internal/infastructur/repo"
	"auth-service/pkg/auth"
	"auth-service/pkg/hash"
	"auth-service/pkg/logger"
	"auth-service/pkg/sender"
	"context"
	"github.com/google/uuid"
	"net"
)

const (
	subjectEmailText = "Изменен ip"
	bodyEmailText    = "Если вы не меняли пароль, то все-го хо-ро-ше-го"
)

type Auth struct {
	userRepo     repo.UserRepo
	sessionRepo  repo.SessionRepo
	transact     repo.Transactor
	tokenManager auth.TokenManager
	hasher       hash.Hasher
	sender       sender.Sender
	logger       logger.Interface
}

func NewAuth(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
	transact repo.Transactor,
	tokenManager auth.TokenManager,
	hasher hash.Hasher,
	sender sender.Sender,
	logger logger.Interface,
) *Auth {
	return &Auth{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		transact:     transact,
		tokenManager: tokenManager,
		hasher:       hasher,
		sender:       sender,
		logger:       logger,
	}
}

func (r *Auth) RefreshOperation(ctx context.Context, dto dto.RefreshOperationRequest) (newAccess string, newRefresh string, err error) {
	claims, err := r.tokenManager.Parse(dto.Access)
	if err != nil {
		r.logger.Warn("token parsing failed: %v", err)
		return "", "", entity.ErrUnauthorized
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		r.logger.Warn("invalid JWT ID: %v", err)
		return "", "", entity.ErrUnauthorized
	}

	session, err := r.sessionRepo.GetByJwtID(ctx, jti)
	if err != nil {
		r.logger.Warn("session not found for jti %s: %v", jti, err)
		return "", "", entity.ErrUnauthorized
	}

	if r.hasher.Compare(session.RefreshTokenHash, dto.Refresh) != nil {
		r.logger.Warn("refresh token mismatch for jti %s", jti)
		return "", "", entity.ErrUnauthorized
	}

	currentIP := dto.IP
	if session.UserIP.String() != currentIP {
		go func(userID uuid.UUID) {
			bgCtx := context.Background()

			emailTo, err := r.userRepo.GetEmailByID(bgCtx, userID)
			if err != nil {
				r.logger.Warn("couldn't get email by userID %s: %v", userID, err)
				return
			}

			err = r.sender.Send(bgCtx, emailTo, subjectEmailText, bodyEmailText)
			if err != nil {
				r.logger.Warn("couldn't send warning email to %s: %v", emailTo, err)
			}
		}(session.UserID)
	}

	newJWT, err := r.tokenManager.GenerateAccess(session.UserID, net.ParseIP(currentIP))
	if err != nil {
		r.logger.Error("failed to generate new access token: %v", err)
		return "", "", entity.ErrGenerateAccess
	}

	newRefresh, expiredRefresh := r.tokenManager.GenerateRefresh(session.ID)
	newRefreshHash, err := r.hasher.Generate(newRefresh)
	if err != nil {
		r.logger.Error("failed to hash new refresh token: %v", err)
		return "", "", entity.ErrGenerateHash
	}

	claimsNew, err := r.tokenManager.Parse(newJWT)
	if err != nil {
		r.logger.Error("failed to parse new JWT: %v", err)
		return "", "", entity.ErrBadParse
	}

	jtiNew, err := uuid.Parse(claimsNew.ID)
	userID, err2 := uuid.Parse(claimsNew.Subject)
	if err != nil || err2 != nil {
		r.logger.Error("invalid new JWT claims: %v | %v", err, err2)
		return "", "", entity.ErrBadParse
	}

	err = r.transact.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := r.sessionRepo.Delete(txCtx, session.ID); err != nil {
			r.logger.Error("failed to delete old session: %v", err)
			return err
		}

		return r.sessionRepo.CreateSession(txCtx, entity.Session{
			ID:               uuid.New(),
			RefreshTokenHash: newRefreshHash,
			Jti:              jtiNew,
			UserIP:           claimsNew.ClientIP,
			ExpiredAt:        expiredRefresh,
			CreatedAt:        claimsNew.IssuedAt.UTC(),
			UserID:           userID,
		})
	})
	if err != nil {
		r.logger.Error("transaction failed: %v", err)
		return "", "", entity.ErrTransaction
	}

	return newJWT, newRefresh, nil
}

func (r *Auth) CreatePairTokens(ctx context.Context, request dto.AccessRefreshTokensRequest) (access string, refresh string, err error) {
	currentIP := net.ParseIP(request.IP)
	if currentIP == nil {
		r.logger.Warn("invalid IP provided: %s", request.IP)
		return "", "", entity.ErrIPNotFound
	}

	accessToken, err := r.tokenManager.GenerateAccess(request.UserID, currentIP)
	if err != nil {
		r.logger.Error("failed to generate access token for user %s: %v", request.UserID, err)
		return "", "", entity.ErrGenerateAccess
	}

	claims, err := r.tokenManager.Parse(accessToken)
	if err != nil {
		r.logger.Error("failed to parse access token for user %s: %v", request.UserID, err)
		return "", "", entity.ErrBadParse
	}

	sessionID := uuid.New()
	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		r.logger.Error("failed to parse JTI from token claims: %v", err)
		return "", "", entity.ErrBadParse
	}

	refreshToken, expiredRefresh := r.tokenManager.GenerateRefresh(sessionID)

	refreshTokenHash, err := r.hasher.Generate(refreshToken)
	if err != nil {
		r.logger.Error("failed to hash refresh token: %v", err)
		return "", "", entity.ErrGenerateHash
	}

	sessionNew := entity.Session{
		ID:               uuid.New(),
		RefreshTokenHash: refreshTokenHash,
		Jti:              jti,
		UserIP:           currentIP,
		ExpiredAt:        expiredRefresh,
		CreatedAt:        claims.IssuedAt.UTC(),
		UserID:           request.UserID,
	}

	err = r.sessionRepo.CreateSession(ctx, sessionNew)
	if err != nil {
		r.logger.Error("failed to create session for user %s: %v", request.UserID, err)
		return "", "", entity.ErrCreateSession
	}

	return accessToken, refreshToken, nil
}
