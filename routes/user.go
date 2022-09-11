package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nurivansyah/weekend-api/controllers"
	"github.com/nurivansyah/weekend-api/middleware"
)

func UserRoutes(router *gin.Engine) {
	r := router.Group("/")
	r.Use(middleware.RequireAuth, middleware.RequireAdmin)
	{
		r.POST("/users", controllers.CreateUser)
		r.GET("/users", controllers.UserIndex)
		r.GET("/users/:user_id", controllers.ShowUser)
		r.PUT("/users/:user_id", controllers.UpdateUser)
		r.DELETE("/users/:user_id", controllers.DeleteUser)
	}
}
