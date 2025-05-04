package dto

import "github.com/google/uuid"

type RefreshOperationRequest struct {
	Refresh string `json:"refresh" binding:"required"`
	Access  string `json:"access" binding:"required"`
	IP      string
}

type RefreshOperationResponse struct {
	Refresh string `json:"refresh" binding:"required"`
	Access  string `json:"access" binding:"required"`
}

type AccessRefreshTokensRequest struct {
	UserID uuid.UUID
	IP     string
}

type AccessRefreshTokensResponse struct {
	Refresh string `json:"refresh" binding:"required"`
	Access  string `json:"access" binding:"required"`
}
