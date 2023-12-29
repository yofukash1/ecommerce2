package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yofukashi/e-commerce/internal/usecase"
	"github.com/yofukashi/e-commerce/pkg/logging"
)

func NewRouter(handler *gin.Engine, e usecase.EcommerceUseCaseI, l *logging.Logger) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newEcommerceRoutes(h, e, l)
	}
}
