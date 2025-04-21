package app

import (
	"net/http"

	"smolink/internal/config"
	"smolink/internal/controller"
	"smolink/internal/repository"
	"smolink/internal/routes"
	"smolink/internal/service"
	"smolink/pkg/database"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router        *gin.Engine
	PGRepo        *repository.PostgresRepository
	RedisRepo     *repository.RedisRepository
	URLService    *service.URLService
	URLController *controller.URLController
	DBCloser      func() error
}

func NewApp(cfg *config.Config, includeRootRoutes bool) (*App, error) {
	pgDB, err := database.NewPostgresDB(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	redisClient, err := database.NewRedisDB(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		return nil, err
	}

	pgRepo := repository.NewPostgresRepository(pgDB.Pool)
	redisRepo := repository.NewRedisRepository(redisClient.Client)
	urlService := service.NewURLService(pgRepo, redisRepo)
	urlController := controller.NewURLController(urlService)

	router := gin.New()

	routes.SetupRoutes(router, urlController)

	if includeRootRoutes {
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Ah, you don reach home. Welcome to the smolink service. We dey for you!",
			})
		})

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Omo, this route no dey exist o! You don miss road. Go back jare!",
			})
		})

		router.NoMethod(func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"message": "Abeg, no try that method here. Na wrong move!",
			})
		})
	}

	return &App{
		Router:        router,
		PGRepo:        pgRepo,
		RedisRepo:     redisRepo,
		URLService:    urlService,
		URLController: urlController,
		DBCloser: func() error {
			pgDB.Close()
			return nil
		},
	}, nil
}
