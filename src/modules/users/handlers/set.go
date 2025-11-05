package handlers

import (
	"bismillah/src/modules/users/dto"
	"bismillah/src/modules/users/model"
	response "bismillah/src/utils/response_helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  *model.UserService
	validate *validator.Validate
}

// Constructor
func NewUserHandler(s *model.UserService) *UserHandler {
	validate := validator.New()
	dto.RegisterValidators(validate) // register custom validator

	return &UserHandler{
		service:  s,
		validate: validate,
	}
}

// --------------------
// RegisterUser handler
// --------------------
func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var payload dto.RegisterUserDTO

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	// Validasi input dengan custom validator
	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(ctx, err.Error(), "Validation failed")
		return
	}

	// Buat user baru
	id, err := h.service.CreateUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		// Tangani error berdasarkan isi pesan (tanpa exceptions)
		switch {
		case err.Error() == "email already registered":
			response.Conflict(ctx, err.Error(), nil)
		default:
			response.InternalServerError(ctx, "Failed to create user", err.Error())
		}
		return
	}

	// Berhasil dibuat
	response.Created(ctx, map[string]string{"id": id}, "User created successfully", nil)
}

// --------------------
// LoginUser handler
// --------------------
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var payload dto.LoginUserDTO

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	// Validasi input login
	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(ctx, err.Error(), "Validation failed")
		return
	}

	// Login user
	token, err := h.service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		switch {
		case err.Error() == "invalid credentials":
			response.Unauthorized(ctx, err.Error(), nil)
		case err.Error() == "user not found":
			response.NotFound(ctx, err.Error(), nil)
		default:
			response.InternalServerError(ctx, "Login failed", err.Error())
		}
		return
	}

	// Sukses login
	response.Success(ctx, map[string]string{"token": token}, "Login successful", 200, nil)
}
