package handler

import (
	"net/http"

	"bismillah/src/api/users/dto"
	"bismillah/src/exceptions"
	"bismillah/src/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(s *service.UserService) *UserHandler {
	validate := validator.New()
	dto.RegisterValidators(validate) // register custom validator

	return &UserHandler{
		service:  s,
		validate: validate,
	}
}

// RegisterUser handler
func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var payload dto.RegisterUserDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi custom password & rules lain
	if err := h.validate.Struct(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		if _, ok := err.(*exceptions.ClientError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"id":      id,
	})
}

// LoginUser handler
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var payload dto.LoginUserDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		if _, ok := err.(*exceptions.AuthenticationError); ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// GetAllUsers handler
func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}
