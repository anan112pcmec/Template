package authway

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func Login(db *gorm.DB, nama, password string) map[string]interface{} {
	var user models.User

	result := db.Where("nama = ? AND password = ?", nama, password).First(&user)

	if result.Error != nil {
		fmt.Println("Login gagal:", result.Error)
		return map[string]interface{}{
			"status":  "false",
			"message": "Nama atau password salah",
		}
	}

	fmt.Println("Login berhasil:", user)

	return map[string]interface{}{
		"status":     "true",
		"message":    "Login berhasil",
		"Nama":       user.Nama,
		"Password":   user.Password,
		"ID":         user.ID,
		"Favorit":    user.GenreDisukai,
		"KreditSkor": fmt.Sprintf("%d", user.KreditSkor),
	}
}
