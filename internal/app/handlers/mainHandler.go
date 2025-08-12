package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

)

func MainPageHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "main.html", nil)

}

func SignUpPageHandler(c *gin.Context){
	c.HTML(http.StatusOK, "signup.html", nil)
}