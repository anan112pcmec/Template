package authway

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID         string `gorm:"primaryKey"`
	Nama       string
	Password   string
	KreditSkor int8 `gorm:"column:kreditskor"`
}

func Login(db *gorm.DB, nama, password string) map[string]string {
	var user User

	result := db.Where("nama = ? AND password = ?", nama, password).First(&user)

	if result.Error != nil {
		fmt.Println("Login gagal:", result.Error)
		return map[string]string{
			"status":  "false",
			"message": "Nama atau password salah",
		}
	}

	fmt.Println("Login berhasil:", user)

	return map[string]string{
		"status":     "true",
		"message":    "Login berhasil",
		"Nama":       user.Nama,
		"Password":   user.Password,
		"ID":         user.ID,
		"KreditSkor": fmt.Sprintf("%d", user.KreditSkor),
	}
}
