package service

import (
	"database/sql"
	"time"
	"os"

	"bismillah/src/exceptions"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO: type UserService struct diibaratkan seperti "class UserService" dalam JavaScript. "db" disini setara dengan this.db di JavaScript. 
type UserService struct {
	db *sql.DB
}

// TODO: NewUserService membuat instance baru dari UserService dengan koneksi database yang diberikan. Instance ini bisa digunakan untuk memanggil method-method seperti CreateUser, LoginUser, GetAllUsers. Ini juga berperan seperti constructor(db), yang nantinya this.db di fungsi UsersService dapat dijalankan.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser service
func (s *UserService) CreateUser(name, email, password string) (string, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists) // TODO: "&" adalah penunjuk ke alamat / memory variable disimpan. Jadi disini maksudnya masukan hasil query ke dalam "exists" dengan alamat memorynya direpresentasikan oleh "&".
	if err != nil {
		return "", err
	}
	if exists {
		return "", &exceptions.ClientError{Message: "Email already registered"}
	}

	// Generate id sendiri pakai ksuid
	id := ksuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert data dan ambil id dari RETURNING
	var userId string
	err = s.db.QueryRow(
		"INSERT INTO users(id, name, email, password) VALUES($1, $2, $3, $4) RETURNING id",
		id, name, email, string(hashedPassword),
	).Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}


// LoginUser service
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