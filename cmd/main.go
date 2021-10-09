package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/online-shop/internal/auth"
	"github.com/online-shop/internal/config"
	"github.com/online-shop/internal/errors"
	"github.com/online-shop/internal/healthcheck"
	"github.com/online-shop/pkg/accesslog"
	"github.com/online-shop/pkg/log"
	"net/http"
	"os"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/go-sql-driver/mysql"
)

var Version = "0.1.0"

func main() {
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	db, err := buildMysqlClient(cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, cfg, db),
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, cfg *config.Config, db *sqlx.DB) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	//authHandler := auth.Handler(cfg.JWTSigningKey)

	//album.RegisterHandlers(rg.Group(""),
	//	album.NewService(album.NewRepository(db, logger), logger),
	//	authHandler, logger,
	//)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	return router
}

func buildMysqlClient(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", cfg.DSN)
	if err != nil {
		return db, err
	}
	return db, nil
}
