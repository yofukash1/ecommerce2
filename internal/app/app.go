package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/yofukashi/e-commerce/internal/config"
	v1 "github.com/yofukashi/e-commerce/internal/controllers/http/v1"
	"github.com/yofukashi/e-commerce/internal/usecase"
	"github.com/yofukashi/e-commerce/internal/usecase/repo"
	"github.com/yofukashi/e-commerce/pkg/httpserver"
	"github.com/yofukashi/e-commerce/pkg/logging"
	"github.com/yofukashi/e-commerce/pkg/postgresql"
)

func Run(ctx context.Context, cfg *config.Config) {
	l := logging.GetLogger()

	l.Info("connecting to db")
	pg, err := postgresql.NewClient(ctx, 10, cfg.Storage)

	if err != nil {
		l.Fatalf("error while connecting to db: %v", err)
	}
	defer pg.Pool.Close()

	l.Info("creating ecommerceUseCase")
	ecommerceUseCase := usecase.New(l, repo.New(pg, l))

	// HTTP Server
	l.Info("starting http server")
	handler := gin.New()
	v1.NewRouter(handler, ecommerceUseCase, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Listen.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("signal: %w", err))

	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
