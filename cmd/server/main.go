package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"bismillah/src/api/users/routes"
)

func main() {
	//TODO: godotenv.Load() berguna untuk membaca file .env di direktori root project.
	//TODO: variable "err" biasa digunakan untuk menangani error yang mungkin terjadi saat memuat function yang bisa saja mengembalikan error. Tujuannya sama seperti try-catch di JS yaitu untuk menangkap error.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//TODO: gin.Default() = bikin router baru (objek pengatur route HTTP) yang sudah siap pakai karena otomatis dikasih logging dan recovery middleware. Logger() → mencatat setiap request, Recovery() → menangani panic agar server tidak mati.
	//TODO: jadi karena basicnya saya ingin ada routes yang nantinya diharapkan dapat di akses, maka saya buatkan dulu "gerbang" nya menggunakan gin.Default().
	//TODO: Variable "r" di sini bertipe *gin.Engine.
	r := gin.Default()

	// Register user routes
	//TODO: routes adalah package, RegisterUserRoutes adalah function yang ada di package routes.
	routes.RegisterUserRoutes(r)

	//TODO: os.Getenv("NAMA_VARIABLE") dipakai untuk: Mengambil nilai dari environment variable yang ada di sistem atau yang sudah dimuat lewat .env.
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running on port", port)
	r.Run(":" + port)
}
