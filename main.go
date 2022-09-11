package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nurivansyah/weekend-api/docs"
	"github.com/nurivansyah/weekend-api/initializers"
	"github.com/nurivansyah/weekend-api/routes"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDb()
}

// @securityDefinitions.apikey bearerAuth
// @in                         header
// @name                       Authorization
func main() {

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = os.Getenv("APP_NAME")
	docs.SwaggerInfo.Description = "This contains minimum API for test purposes."
	docs.SwaggerInfo.Version = "master"
	docs.SwaggerInfo.Host = os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT")
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()

	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from " + os.Getenv("APP_NAME"),
		})
	})

	r.Run()
}
