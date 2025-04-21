package routes

import (
	"net/http"
	"smolink/internal/controller"
	"smolink/pkg/middleware"

	"github.com/gin-gonic/gin"
)

const (
	APIPrefix       = "/api/v1"
	ShortenURLPath  = "/links"
	HealthCheckPath = "/health"
)

func SetupUrlRoutes(router *gin.Engine, urlController *controller.URLController) {
	urlGroup := router.Group(APIPrefix)
	{
		urlGroup.POST(ShortenURLPath, urlController.ShortenURL)
		urlGroup.GET(ShortenURLPath+"/:code", urlController.ResolveURL)
	}
}

func SetupRoutes(router *gin.Engine, urlController *controller.URLController) {
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())

	router.GET(HealthCheckPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	SetupUrlRoutes(router, urlController)
}
