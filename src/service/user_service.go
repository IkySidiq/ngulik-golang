package service

import (
	"database/sql"
	"time"
	"os"

	"bismillah/src/exceptions"
	"bismillah/src/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService seperti class UserService di JS, punya akses ke this.db
type UserService struct {
	db *sql.DB
}

// Constructor untuk UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser service (sekarang pakai BaseModel)
func (s *UserService) CreateUser(name, email, password string) (string, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "", &exceptions.ClientError{Message: "Email already registered"}
	}

	// --- Bagian penting: gunakan BaseModel ---
	userModel := utils.NewBaseModel(s.db, "users")

	id := ksuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Data yang akan diinsert
	userData := map[string]interface{}{
		"id":       id,
		"name":     name,
		"email":    email,
		"password": string(hashedPassword),
	}

	// Gunakan BaseModel.Create
	createdUser, err := userModel.Create(userData)
	if err != nil {
		return "", err
	}

	// Kembalikan id hasil create
	return createdUser["id"].(string), nil
}

// LoginUser service (belum diubah)
func (s *UserService) LoginUser(email, password string) (string, error) {
	var id string
	var hashedPassword string
	err := s.db.QueryRow("SELECT id, password FROM users WHERE email=$1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return "", &exceptions.AuthenticationError{Message: "Invalid email or password"}
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return "", &exceptions.AuthenticationError{Message: "Invalid email or password"}
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString, nil
}

// GetAllUsers service
func (s *UserService) GetAllUsers() ([]map[string]interface{}, error) {
	rows, err := s.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id string
		var name, email string
		rows.Scan(&id, &name, &email)
		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	return users, nil
}