package v1

import (
	"auth-service/internal/controller/http/middleware"
	"auth-service/internal/controller/http/v1/auth"
	"auth-service/internal/usecase"
	"auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	l logger.Interface,
	authUC usecase.Auth,
) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.Logger(l),
	)

	apiV1 := router.Group("/v1")
	{
		auth.NewRouter(
			apiV1,
			authUC,
			l,
		)
	}

	return router
}
