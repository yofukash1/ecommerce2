package main

import (
	"context"

	"github.com/yofukashi/e-commerce/internal/app"
	"github.com/yofukashi/e-commerce/internal/config"
)

func main() {
	// Configuration
	cfg := config.GetConfig()
	ctx := context.Background()

	// Run
	app.Run(ctx, cfg)
}
