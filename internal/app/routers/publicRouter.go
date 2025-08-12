package routers

import (
	"log"

	"github.com/gin-gonic/gin"

	"planify/internal/app/handlers"
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

	r.GET("/", handlers.MainPageHandler)

	r.GET("/signup", handlers.SignUpPageHandler)
	r.POST("/signup", userHandler.SignUpHandler)

	r.GET("/login", handlers.LoginPageHandler)
	r.POST("/login", userHandler.LoginHandler)

	return r

}
