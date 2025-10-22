package routes

import (
	"github.com/gin-gonic/gin"
	"bismillah/src/api/users/handler"
	"bismillah/src/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", handler.RegisterUser)
		userGroup.POST("/login", handler.LoginUser)
		userGroup.GET("/", middleware.JWTAuthMiddleware(), handler.GetAllUsers)
	}
}
