package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurivansyah/weekend-api/dto"
)

func RequireAdmin(c *gin.Context) {
	user, _ := c.Get("user")
	user_data := user.(dto.UserSession)

	if user_data.User_type != "ADMIN" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "Unauthorized Access",
		})
		return
	}

	c.Next()
}
