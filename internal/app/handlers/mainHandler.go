package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

)

func IndexPageHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)

}

func SignUpPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}


func MainPageHandler(c *gin.Context){
	c.HTML(http.StatusOK, "main.html", nil)
}