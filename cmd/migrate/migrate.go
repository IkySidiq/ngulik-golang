package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run migrate.go create|up|down <name-if-create>")
		return
	}

	command := os.Args[1]

	// URL database PostgreSQL
	dbURL := "postgres://developer:BackEndDreamer2025@localhost:5432/bismillah?sslmode=disable"
	migrationsPath := "./migrations"

	var cmd *exec.Cmd

	switch command {
	case "up", "down":
		// jalankan migrate up atau down
		cmd = exec.Command("migrate", "-path", migrationsPath, "-database", dbURL, command)

	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run migrate.go create <migration_name>")
			return
		}
		name := os.Args[2]

		// buat timestamp unik untuk nama file
		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s", timestamp, name)

		// perintah migrate create
		cmd = exec.Command("migrate", "create", "-ext", "sql", "-dir", migrationsPath, fileName)

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage: go run migrate.go create|up|down <name-if-create>")
		return
	}

	// tampilkan output di terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running migration:", err)
	}
}