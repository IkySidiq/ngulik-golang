package service

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"bismillah/src/exceptions"

	_ "github.com/lib/pq"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

// CreateUser service
func CreateUser(name, email, password string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return &exceptions.ClientError{Message: "Email already registered"}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = db.Exec("INSERT INTO users(name, email, password) VALUES($1, $2, $3)", name, email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

// LoginUser service
func LoginUser(email, password string) (string, error) {
	var id int
	var hashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE email=$1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return "", &exceptions.AuthenticationError{Message: "Invalid email or password"}
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return "", &exceptions.AuthenticationError{Message: "Invalid email or password"}
	}

	// generate JWT
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString, nil
}

// GetAllUsers service
func GetAllUsers() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
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
