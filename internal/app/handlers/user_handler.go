package handlers

import (
	"log"
	"net/http"
	"planify/internal/domain/models"
	"planify/internal/domain/repo"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repo.UserRepoInterface
}

func NewUserHandler(userRepo repo.UserRepoInterface) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user internal"})
		return
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), &user); err != nil {
		log.Printf("Error create user:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
