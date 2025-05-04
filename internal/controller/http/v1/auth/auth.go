package auth

import (
	"auth-service/internal/controller/http/dto"
	"auth-service/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (r *Router) Authentification(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.logger.Warn(dto.MsgInvalidParams, err)
		dto.ErrorResponse(c, http.StatusBadRequest, dto.MsgInvalidParams)
		return
	}

	req := dto.AccessRefreshTokensRequest{
		UserID: userID,
		IP:     c.ClientIP(),
	}

	access, refresh, err := r.authUC.CreatePairTokens(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrIPNotFound),
			errors.Is(err, entity.ErrGenerateAccess),
			errors.Is(err, entity.ErrBadParse),
			errors.Is(err, entity.ErrCreateSession):

			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())

		default:
			r.logger.Error("Router - Authentification - r.authUC.CreatePairTokens():", err)

			dto.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.AccessRefreshTokensResponse{
		Refresh: refresh,
		Access:  access,
	}

	c.JSON(http.StatusOK, resp)
}
