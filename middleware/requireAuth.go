package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nurivansyah/weekend-api/dto"
	"github.com/nurivansyah/weekend-api/initializers"
	"github.com/nurivansyah/weekend-api/models"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func RequireAuth(c *gin.Context) {

	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "authorization header is not provided",
		})
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "authorization header format invalid",
		})
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: fmt.Sprintf("authorization type unsupported: %s", authorizationType),
		})
		return
	}

	tokenString := fields[1]

	jwt_string := strings.Split(tokenString, ".")
	if len(jwt_string) != 3 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "jwt format invalid",
		})
		return
	}

	// tokenString, err := c.Cookie("Authorization")
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// }

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}

		return []byte(os.Getenv("TOKEN_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "access token expired",
			})
			return
		}

		if claims["aud"] != os.Getenv("TOKEN_AUDIENCE") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "access token payload invalid",
			})
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "access token payload invalid",
			})
			return
		}

		user_session := dto.UserSession{
			ID:        user.ID,
			Username:  user.Username,
			User_type: user.User_type,
		}

		c.Set("user", user_session)

		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "access token invalid/expired",
		})
		return
	}
}
