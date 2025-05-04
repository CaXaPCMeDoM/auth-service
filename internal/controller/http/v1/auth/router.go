package auth

import (
	"auth-service/internal/usecase"
	"auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authUC usecase.Auth
	logger logger.Interface
}

func NewRouter(
	apiV1Group *gin.RouterGroup,
	authUC usecase.Auth,
	logger logger.Interface,
) *Router {
	au := &Router{
		authUC: authUC,
		logger: logger,
	}

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/tokens", au.Authentification)
		authGroup.POST("/refresh", au.Refresh)
	}

	return au
}
