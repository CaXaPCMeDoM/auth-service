package app

import (
	"auth-service/config"
	v1 "auth-service/internal/controller/http/v1"
	"auth-service/internal/infrastructure/repo/persistence"
	"auth-service/internal/usecase/auth"
	"auth-service/pkg/auth/jwt"
	"auth-service/pkg/hash/bcrypt"
	"auth-service/pkg/httpserver"
	"auth-service/pkg/logger"
	"auth-service/pkg/postgres"
	"auth-service/pkg/sender/email"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// tools
	l := logger.New(cfg.Log.Level)
	tokenManager, err := jwt.New(
		[]byte(cfg.Jwt.SecretKey),
		cfg.Jwt.AccessTokenTTL,
		cfg.Jwt.RefreshTokenTTL,
	)
	hasher := bcrypt.NewBcryptHasher(cfg.Security.PasswordCost)
	sender, err := email.New(*cfg)

	// repo
	pg, err := postgres.New(cfg.Pg.URL, postgres.MaxPoolSize(cfg.Pg.PoolMax))

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	defer pg.Close()

	userRepo := persistence.NewUserRepo(pg)
	sessionRepo := persistence.NewSessionRepo(pg)

	// usecase
	authUC := auth.NewAuth(
		userRepo,
		sessionRepo,
		pg,
		tokenManager,
		hasher,
		sender,
		l,
	)

	//controllers
	router := v1.NewRouter(
		l,
		authUC,
	)

	//HTTp server
	server := httpserver.New(
		cfg,
		router,
		httpserver.Mode(cfg.HTTP.Mode),
	)

	server.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
