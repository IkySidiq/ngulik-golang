package routes

import (
	"github.com/gin-gonic/gin"
	"bismillah/src/modules/users/handlers"
	"bismillah/src/modules/users/model"
	"bismillah/src/utils"
	"bismillah/src/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	db, _ := utils.NewDB() // ambil koneksi DB
	userService := model.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
		userGroup.GET("/", middleware.JWTAuthMiddleware(), userHandler.GetAllUsers)
	}
}
