package routes

import (
	"github.com/gin-gonic/gin"
	"bismillah/src/api/users/handler"
	"bismillah/src/service"
	"bismillah/src/utils"
	"bismillah/src/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	db, _ := utils.NewDB() // ambil koneksi DB
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
		userGroup.GET("/", middleware.JWTAuthMiddleware(), userHandler.GetAllUsers)
	}
}
