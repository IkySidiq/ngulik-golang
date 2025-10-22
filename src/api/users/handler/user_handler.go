package handler

import (
	"net/http"

	"bismillah/src/service"
	"bismillah/src/exceptions"

	"github.com/gin-gonic/gin"
)


// RegisterUser handler
func RegisterUser(c *gin.Context) {
	var payload struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	//* if <short variable declaration>; <condition>
	if err := c.ShouldBindJSON(&payload); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		if _, ok := err.(*exceptions.ClientError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// LoginUser handler
func LoginUser(c *gin.Context) {
	var payload struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		if _, ok := err.(*exceptions.AuthenticationError); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetAllUsers handler
func GetAllUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
