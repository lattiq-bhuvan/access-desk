package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/lattiq/foundry/config"
	"github.com/lattiq/foundry/o11y/logging"
	"github.com/lattiq/foundry/service"
	_http "github.com/lattiq/foundry/service/http"
	_gin "github.com/lattiq/foundry/service/http/gin"
	"github.com/lattiq/foundry/validator"
)

type Config struct {
	Logging logging.Config `json:"logging"`
	Server  _http.Config   `json:"server"`
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
	serverName := "access-desk"
	svc := service.New(serverName).WithShutdownTimeout(5 * time.Second)

	router := _gin.Router(&cfg.Server)
	httpServer := _http.NewServer(&cfg.Server, router)
	svc.AddServer("http", httpServer, cfg.Server.Address())

	if err := svc.Run(context.Background()); err != nil {
		slog.Error("service stopped with error", "error", err)
		os.Exit(1)
	}
}
