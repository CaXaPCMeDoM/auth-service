package auth

import (
	"auth-service/internal/controller/http/dto"
	"auth-service/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Router) Refresh(c *gin.Context) {
	var req dto.RefreshOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.logger.Warn(dto.MsgInvalidRequestBody, err)
		dto.ErrorResponse(c, http.StatusBadRequest, dto.MsgInvalidRequestBody)
		return
	}

	req.IP = c.ClientIP()

	access, refresh, err := r.authUC.RefreshOperation(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrTransaction):
			r.logger.Warn("transaction error: %v", err)
			dto.ErrorResponse(c, http.StatusLocked, dto.MsgParallelAccess)

		case errors.Is(err, entity.ErrUnauthorized):
			r.logger.Warn("unauthorized access: %v", err)
			dto.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")

		case errors.Is(err, entity.ErrBadParse),
			errors.Is(err, entity.ErrSessionNotFound),
			errors.Is(err, entity.ErrInvalidRefreshToken),
			errors.Is(err, entity.ErrGenerateHash),
			errors.Is(err, entity.ErrIPNotFound):
			r.logger.Warn("bad request: %v", err)
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())

		default:
			r.logger.Error("unexpected error: %v", err)
			dto.ErrorResponse(c, http.StatusInternalServerError, dto.MsgInternalError)
		}
		return
	}

	resp := dto.RefreshOperationResponse{
		Access:  access,
		Refresh: refresh,
	}

	c.JSON(http.StatusOK, resp)
}
