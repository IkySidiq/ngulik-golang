package model

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"bismillah/src/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

// Constructor
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser service
func (s *UserService) CreateUser(name, email, password string) (string, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("email already registered")
	}

	userModel := utils.NewBaseModel(s.db, "users")

	id := ksuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	userData := map[string]interface{}{
		"id":       id,
		"name":     name,
		"email":    email,
		"password": string(hashedPassword),
	}

	createdUser, err := userModel.Create(userData)
	if err != nil {
		return "", err
	}

	return createdUser["id"].(string), nil
}

// LoginUser service
func (s *UserService) LoginUser(email, password string) (string, error) {
	var id, hashedPassword string
	err := s.db.QueryRow("SELECT id, password FROM users WHERE email=$1", email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return "", errors.New("invalid email or password")
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// GetAllUsers service
func (s *UserService) GetAllUsers(page, limit int, search string) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, name, email 
		FROM users 
		WHERE name ILIKE $1 
		ORDER BY name 
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return nil, 0, err
		}
		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	var total int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE name ILIKE $1", "%"+search+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
