package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurivansyah/weekend-api/dto"
	"github.com/nurivansyah/weekend-api/initializers"
	"github.com/nurivansyah/weekend-api/models"
	"golang.org/x/crypto/bcrypt"
)

// Paths Information

// CreateUser godoc
// @Security    bearerAuth
// @Summary     Create User
// @Description Create a new user (require authenticated and admin level to call this resource)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       user body     dto.CreateUserRequest true "Create user"
// @Success     201  {object} models.User
// @Failure     400  {object} dto.ErrorResponse
// @Failure     401  {object} dto.MessageResponse
// @Failure     500  {object} dto.ErrorResponse
// @Router      /users [post]
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to read request",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to hash password",
		})
		return
	}

	user := models.User{Username: req.Username, Password: string(hash), User_type: req.User_type}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result.Error.Error(),
			Message: "failed to create request",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UserIndex godoc
// @Security    bearerAuth
// @Summary     Retrieve list of User
// @Description Get all the existing users (require authenticated and admin level to call this resource)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Success     200 {object} []models.User
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.MessageResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /users [get]
func UserIndex(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result.Error.Error(),
			Message: "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ShowUser godoc
// @Security    bearerAuth
// @Summary     Retrieve User
// @Description Show a single user (require authenticated and admin level to call this resource)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       id  path     int true "User ID"
// @Success     200 {object} models.User
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /users/{id} [get]
func ShowUser(c *gin.Context) {
	user_id := c.Param("user_id")

	var user models.User
	result := initializers.DB.First(&user, user_id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result.Error.Error(),
			Message: "failed to fetch a user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Security    bearerAuth
// @Summary     Update User
// @Description Update a single user (require authenticated and admin level to call this resource)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       id   path     int                   true "User ID"
// @Param       user body     dto.UpdateUserRequest true "Update user"
// @Success     200  {object} models.User
// @Failure     400  {object} dto.ErrorResponse
// @Failure     401  {object} dto.MessageResponse
// @Failure     500  {object} dto.ErrorResponse
// @Router      /users/{id} [put]
func UpdateUser(c *gin.Context) {
	user_id := c.Param("user_id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to read request",
		})
		return
	}

	var user models.User
	result_query := initializers.DB.First(&user, user_id)
	if result_query.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result_query.Error.Error(),
			Message: "failed to fetch a user",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   err.Error(),
			Message: "failed to hash password",
		})
		return
	}

	result_update := initializers.DB.Model(&user).Updates(models.User{
		Username:  req.Username,
		Password:  string(hash),
		User_type: req.User_type,
	})
	if result_update.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result_update.Error.Error(),
			Message: "failed to update a user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Security    bearerAuth
// @Summary     Delete User
// @Description Delete a single user (require authenticated and admin level to call this resource)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       id  path     int true "User ID"
// @Success     204 {object} interface{}
// @Failure     401 {object} dto.MessageResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	user_id := c.Param("user_id")

	result := initializers.DB.Delete(&models.User{}, user_id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error:   result.Error.Error(),
			Message: "failed to delete a user",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
