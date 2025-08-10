package routers

import (
	"planify/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPubRouter() *gin.Engine {

	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.GET("/", handlers.MainPageHandler)
	return r

}
