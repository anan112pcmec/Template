package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

var jwtSecret []byte

func init() {
	secretFromEnv := os.Getenv("JWT_SECRET")
	if secretFromEnv != "" {
		jwtSecret = []byte(secretFromEnv)
		log.Println("JWT secret loaded from environment variable")
	} else {
		// Generate 32 bytes random key untuk HS256
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			log.Fatalf("Gagal generate jwt secret: %v", err)
		}
		jwtSecret = key
		log.Printf("JWT Secret auto-generated (base64): %s\n", base64.StdEncoding.EncodeToString(jwtSecret))
	}
}

// HashPassword mengenkripsi password sebelum disimpan
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword membandingkan password input dan hash di DB
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

// GenerateJWT membuat token JWT untuk user yang login
func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// LoginUser mengecek kredensial dan mengembalikan token jika sukses
func LoginUser(db *gorm.DB, usernameOrEmail, password string) (string, error) {
	var user models.User
	err := db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	return GenerateJWT(user)
}
