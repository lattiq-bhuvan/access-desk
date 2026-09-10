package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	govalidator "github.com/go-playground/validator/v10"
	fjwt "github.com/lattiq/foundry/auth/jwt"
	"github.com/lattiq/foundry/config"
	"github.com/lattiq/foundry/database/sql"
	"github.com/lattiq/foundry/o11y/logging"
	"github.com/lattiq/foundry/service"
	_http "github.com/lattiq/foundry/service/http"
	_gin "github.com/lattiq/foundry/service/http/gin"
	fmw "github.com/lattiq/foundry/service/http/middleware"
	"github.com/lattiq/foundry/validator"

	"github.com/lattiq-bhuvan/access-desk/internal/auth"
	"github.com/lattiq-bhuvan/access-desk/internal/handler"
	_service "github.com/lattiq-bhuvan/access-desk/internal/service"
	"github.com/lattiq-bhuvan/access-desk/internal/store"
)

type Config struct {
	Logging  logging.Config `json:"logging"`
	Server   _http.Config   `json:"server"`
	Database sql.Config     `json:"database"`
	JWT      fjwt.Config    `json:"jwt"`
}

func main() {
	// create a new validator instance and register custom validations based on rules in the validator package
	v := govalidator.New()
	if err := validator.RegisterCustomValidations(v, validator.CustomValidations...); err != nil {
		slog.Error("failed to register custom validations", "error", err)
		os.Exit(1)
	}

	var cfg Config
	if err := config.Load(&cfg,
		config.WithValidator(v),
		config.WithEnvPrefix("ACCESSDESK_"),
		config.WithEnvDelimiter("__"),
	); err != nil {
		slog.Error("failed to load and validate config", "error", err)
		os.Exit(1)
	}

	cfg.Logging.Configure()
	serverName := "accessdesk"
	svc := service.New(serverName).WithShutdownTimeout(5 * time.Second)

	st, err := store.New(context.Background(), &cfg.Database)
	if err != nil {
		slog.Error("failed to init store", "error", err)
		os.Exit(1)
	}

	reqSvc := _service.NewRequestService(st)
	reqH := handler.NewRequestHandler(reqSvc)

	accessTTL, _ := time.ParseDuration(cfg.JWT.AccessTokenTTL)
	tokenSvc := fjwt.NewTokenService(cfg.JWT)

	authSvc := auth.NewAuthService(st, tokenSvc, accessTTL)
	authH := handler.NewAuthHandler(authSvc)

	dsH := handler.NewDatasetHandler(st)

	jwtMW := fmw.NewJWTMiddleware(fmw.JWTMiddlewareConfig{
		TokenService:  tokenSvc,
		ClaimsFactory: auth.ClaimsFactory,
	}).Handler()

	router := _gin.Router(&cfg.Server)

	// public routes
	authRouter := router.Group("/auth/v1")
	{
		authRouter.POST("/login", authH.Login)
		authRouter.POST("/refresh", authH.Refresh)
		authRouter.POST("/logout", authH.Logout)
	}

	// protected routes
	v1 := router.Group("/v1", jwtMW)
	{
		v1.GET("/datasets", dsH.List)
		v1.GET("/users/me", authH.Me)
		v1.PATCH("/users/me/settings", authH.UpdateSettings)
		v1.POST("/requests", reqH.Create)
		v1.GET("/requests", reqH.List)
		v1.GET("/requests/:id", reqH.Get)
		v1.POST("/requests/:id/decision", reqH.Decide)
	}

	httpServer := _http.NewServer(&cfg.Server, router)
	svc.AddServer("http", httpServer, cfg.Server.Address())

	if err := svc.Run(context.Background()); err != nil {
		slog.Error("service stopped with error", "error", err)
		os.Exit(1)
	}
}
