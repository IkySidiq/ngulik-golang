package handlers

import (
	"bismillah/src/modules/users/dto"
	"bismillah/src/modules/users/model"
	"bismillah/src/exceptions"
	 response "bismillah/src/utils/response_helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  *model.UserService
	validate *validator.Validate
}

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

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	// Validasi custom password & rules lain
	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(ctx, err.Error(), "Validation failed")
		return
	}

	id, err := h.service.CreateUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		switch err.(type) {
		case *exceptions.ClientError:
			response.BadRequest(ctx, err.Error(), nil)
		default:
			response.InternalServerError(ctx, "Internal server error", nil)
		}
		return
	}

	response.Created(ctx, map[string]string{"id": id}, "User created successfully", nil)
}

// --------------------
// LoginUser handler
// --------------------
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var payload dto.LoginUserDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(ctx, err.Error(), "Validation failed")
		return
	}

	token, err := h.service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		switch err.(type) {
		case *exceptions.AuthenticationError:
			response.Unauthorized(ctx, err.Error(), nil)
		default:
			response.InternalServerError(ctx, "Internal server error", nil)
		}
		return
	}

	response.Success(ctx, map[string]string{"token": token}, "Login successful", 200, nil)
}
