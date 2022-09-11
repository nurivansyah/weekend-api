package main

import (
	"os"

	"github.com/nurivansyah/weekend-api/initializers"
	"github.com/nurivansyah/weekend-api/models"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDb()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})

	cau := os.Getenv("CREATE_ADMIN_USER")
	if cau == "true" {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), 10)
		user := models.User{Username: "admin", Password: string(hash), User_type: "ADMIN"}
		initializers.DB.Create(&user)
	}

}

func CreateAdminUser() {

}
