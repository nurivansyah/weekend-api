package controllers

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
	"golang.org/x/crypto/bcrypt"
)

// Paths Information

// SignUp godoc
// @Summary     SignUp
// @Description Register a user to access API
// @Tags        Auth
// @ID          Signup
// @Accept      json
// @Produce     json
// @Param       signup body     dto.SignUpRequest true "User SignUp"
// @Success     201    {object} interface{}
// @Failure     400    {object} dto.ErrorResponse
// @Router      /auth/signup [post]
func Signup(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to read request",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to hash password",
		})
		return
	}

	user := models.User{Username: req.Username, Password: string(hash), User_type: "USER"}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   result.Error.Error(),
			Message: "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// LogIn godoc
// @Summary     Login
// @Description Authenticates a user and provides access_token and refresh_token to Authorize API calls
// @Tags        Auth
// @ID          Login
// @Accept      json
// @Produce     json
// @Param       login body     dto.LoginRequest true "User Login"
// @Success     200   {object} dto.JWT
// @Failure     400   {object} dto.ErrorResponse
// @Router      /auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to read request",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "username = ?", req.Username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Message: "invalid username or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to hash password",
		})
		return
	}

	access_token_duration, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": os.Getenv("TOKEN_AUDIENCE"),
		"sub": user.ID,
		"exp": time.Now().Add(access_token_duration).Unix(),
	})

	accessTokenString, err := access_token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to create token",
		})
		return
	}

	refresh_token_duration, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": os.Getenv("TOKEN_AUDIENCE") + "/rt",
		"sub": user.ID,
		"exp": time.Now().Add(refresh_token_duration).Unix(),
	})

	refreshTokenString, err := refresh_token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to create token",
		})
		return
	}

	initializers.DB.Model(&user).Updates(models.User{
		Refresh_token: refreshTokenString,
	})

	c.JSON(http.StatusOK, &dto.JWT{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 60*10, "", "", false, true)
	// c.JSON(http.StatusOK, gin.H{})
}

// RenewToken godoc
// @Summary     Renew Token
// @Description Authenticates refresh_token and generates new access_token and refresh_token to Authorize API calls
// @Tags        Auth
// @ID          RenewToken
// @Accept      json
// @Produce     json
// @Param       refresh_token body     dto.RenewTokenRequest true "Renew Access Token"
// @Success     200           {object} dto.JWT
// @Failure     400           {object} dto.ErrorResponse
// @Failure     401           {object} dto.MessageResponse
// @Router      /auth/refresh [post]
func RenewToken(c *gin.Context) {
	var req dto.RenewTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to read request",
		})
		return
	}

	jwt_string := strings.Split(req.RefreshToken, ".")
	if len(jwt_string) != 3 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "refresh token format invalid",
		})
		return
	}

	token, _ := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}

		return []byte(os.Getenv("TOKEN_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "refresh token expired",
			})
			return
		}

		if claims["aud"] != os.Getenv("TOKEN_AUDIENCE")+"/rt" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "invalid refresh token payload",
			})
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "invalid refresh token payload",
			})
			return
		}

		if user.Refresh_token != req.RefreshToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
				Message: "mismatch refresh token",
			})
			return
		}

		access_token_duration, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
		access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"aud": os.Getenv("TOKEN_AUDIENCE"),
			"sub": user.ID,
			"exp": time.Now().Add(access_token_duration).Unix(),
		})

		accessTokenString, err := access_token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
				Error:   err.Error(),
				Message: "failed to create token",
			})
			return
		}

		refresh_token_duration, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
		refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"aud": os.Getenv("TOKEN_AUDIENCE") + "/rt",
			"sub": user.ID,
			"exp": time.Now().Add(refresh_token_duration).Unix(),
		})

		refreshTokenString, err := refresh_token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
				Error:   err.Error(),
				Message: "failed to create token",
			})
			return
		}

		initializers.DB.Model(&user).Updates(models.User{
			Refresh_token: refreshTokenString,
		})

		c.JSON(http.StatusOK, &dto.JWT{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		})

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.MessageResponse{
			Message: "refresh token invalid/expired",
		})
		return
	}
}

// Profile godoc
// @Security    bearerAuth
// @Summary     Retrieve Profile Data
// @Description Get authenticated user data
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.UserProfile
// @Failure     401 {object} dto.MessageResponse
// @Router      /profile [get]
func Profile(c *gin.Context) {
	user, _ := c.Get("user")
	user_data := user.(dto.UserSession)

	c.JSON(http.StatusOK, &dto.UserProfile{
		User: user_data,
	})
}
