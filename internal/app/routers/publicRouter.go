package routers

import (
	"planify/internal/app/handlers"
	"planify/internal/domain/infrastructure/database"
	"planify/internal/domain/repo"

	"github.com/gin-gonic/gin"
)

func SetupPubRouter() *gin.Engine {

	userRepo := repo.NewUserRepo(database.DB)

	userHandler := handlers.NewUserHandler(userRepo)

	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.GET("/", handlers.MainPageHandler)

	r.POST("/user", userHandler.CreateUser)
	return r

}
