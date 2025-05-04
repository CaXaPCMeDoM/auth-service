package usecase

import (
	"auth-service/internal/controller/http/dto"
	"context"
)

type (
	Auth interface {
		RefreshOperation(ctx context.Context, dto dto.RefreshOperationRequest) (newAccess string, newRefresh string, err error)
		CreatePairTokens(ctx context.Context, request dto.AccessRefreshTokensRequest) (access string, refresh string, err error)
	}
)
