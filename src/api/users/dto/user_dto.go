package dto

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterUserDTO menerima data registrasi user
type RegisterUserDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,no_space"` // custom validator no_space
	Password string `json:"password" validate:"required,password"`     // custom validator password
}

// LoginUserDTO untuk login
type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email,no_space"`
	Password string `json:"password" validate:"required"`
}

// Custom password validator
func passwordValidator(fl validator.FieldLevel) bool {
	pass := fl.Field().String()

	// minimal 6 karakter
	if len(pass) < 6 {
		return false
	}

	// tidak boleh ada spasi
	if strings.Contains(pass, " ") {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, c := range pass {
		switch {
		case 'a' <= c && c <= 'z':
			hasLower = true
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case '0' <= c && c <= '9':
			hasDigit = true
		}
	}

	return hasLower && hasUpper && hasDigit
}

// Custom validator untuk memastikan tidak ada spasi
func noSpaceValidator(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	// trim spaces di awal/akhir
	str = strings.TrimSpace(str)
	// cek apakah ada spasi di tengah
	return !strings.Contains(str, " ")
}

// Fungsi register validator custom
func RegisterValidators(v *validator.Validate) {
	v.RegisterValidation("password", passwordValidator)
	v.RegisterValidation("no_space", noSpaceValidator)
}
