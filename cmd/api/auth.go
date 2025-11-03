package main

import (
	"log"
	"net/http"
	"rest-api-go-gin/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

func (app *application) registerUser(c *gin.Context) {
	var payload registerRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt the password"})
		return
	}

	payload.Password = string(hashedPassword)

	user := database.User{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
	}

	if err := app.models.Users.Insert(&user); err != nil {
		log.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create a user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
