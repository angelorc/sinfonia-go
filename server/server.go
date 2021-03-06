package server

import (
	"context"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/utility"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"os/signal"
	"time"
)

func Start(cfg config.Config) {
	// Setup & configure server
	// more info -> https://echo.labstack.com/
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: utility.LoggerFormat()}))
	e.HideBanner = true

	// Load routes from graphql
	InitGraphql(cfg, e)

	// Load routes from rest
	InitRest(e)

	// Parse server routes
	go func() {
		if err := e.Start(cfg.GraphQL.Address + ":" + cfg.GraphQL.Port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Stop server gracefully
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
