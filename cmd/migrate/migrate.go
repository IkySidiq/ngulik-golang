package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run migrate.go up|down")
		return
	}

	command := os.Args[1]

	dbURL := "postgres://developer:BackEndDreamer2025@localhost:5432/bismillah?sslmode=disable"
	migrationsPath := "./migrations"

	cmd := exec.Command("migrate", "-path", migrationsPath, "-database", dbURL, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running migration:", err)
	}
}
