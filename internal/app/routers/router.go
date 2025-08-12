package routers

import (
	"log"

	"github.com/gin-gonic/gin"

	"planify/internal/app/handlers"
	"planify/internal/app/middleware"
	"planify/internal/domain/config"
	"planify/internal/domain/infrastructure/database"
	"planify/internal/domain/repo"
)

func SetupPubRouter() *gin.Engine {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	userRepo := repo.NewUserRepo(database.DB)

	userHandler := handlers.NewUserHandler(userRepo, cfg)

	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.Static("static", "./static")

	r.GET("/", handlers.IndexPageHandler)

	r.GET("/signup", handlers.SignUpPageHandler)
	r.POST("/signup", userHandler.SignUpHandler)

	r.GET("/login", handlers.LoginPageHandler)
	r.POST("/login", userHandler.LoginHandler)

	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		// Проверяем токен в URL только для GET запросов
		if c.Request.Method == "GET" {
			if token := c.Query("token"); token != "" {
				c.Request.Header.Set("Authorization", "Bearer "+token)
			}
		}
		c.Next()
	})
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/main", handlers.MainPageHandler)

	}

	return r

}
