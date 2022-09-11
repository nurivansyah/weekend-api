package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nurivansyah/weekend-api/controllers"
	"github.com/nurivansyah/weekend-api/middleware"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/refresh", controllers.RenewToken)
	r.GET("/profile", middleware.RequireAuth, controllers.Profile)
}
