package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"planify/internal/domain/config"
	"planify/internal/domain/models"
	"planify/internal/domain/repo"
)

type UserHandler struct {
	userRepo repo.UserRepoInterface
	cfg      *config.ConfigDB
}

func NewUserHandler(userRepo repo.UserRepoInterface, cfg *config.ConfigDB) *UserHandler {
	return &UserHandler{userRepo: userRepo, cfg: cfg}
}

func (h *UserHandler) SignUpHandler(c *gin.Context) {
	var requers struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if requers.Username == "" || requers.Email == "" || requers.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password are required"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(requers.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	user := models.Users{
		Username:  requers.Username,
		Email:     requers.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), &user); err != nil {
		log.Printf("Error create user:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "user create succesfully",
		"user":    user,
	})
}
