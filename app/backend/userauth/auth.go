package userauth

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

)

// User model minimal, sesuaikan dengan project kamu
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Role     string
}

// HashPassword mengenkripsi password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword mencocokkan password plaintext dengan hash
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

// LoginUser melakukan autentikasi dan mengembalikan token JWT dan info lainnya
func LoginUser(db *gorm.DB, usernameOrEmail, password string) map[string]interface{} {
	start := time.Now()

	var user User
	err := db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": "User tidak ditemukan",
		}
	}

	if !CheckPassword(password, user.Password) {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": "Password salah",
		}
	}

	// Generate JWT token via fungsi dari jwt.go
	claims := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	}
	token, err := CreateJWT(claims, 24*time.Hour)
	if err != nil {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": fmt.Sprintf("Gagal membuat token JWT: %v", err),
		}
	}

	elapsed := time.Since(start).Seconds() * 1000

	return map[string]interface{}{
		"Status":                  "Berhasil",
		"Keterangan":              "Autentikasi Teruji Valid",
		"Waktu_Pengautentikasian": fmt.Sprintf("%.2f milisecond", elapsed),
		"Token":                   token,
		"iduser":                  user.ID,
	}
}

// RegisterUser membuat user baru dan mengembalikan status proses
func RegisterUser(db *gorm.DB, username, email, password string) map[string]interface{} {
	start := time.Now()

	// Cek apakah user/email sudah ada
	var existing User
	err := db.Where("username = ? OR email = ?", username, email).First(&existing).Error
	if err == nil {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": "Username atau Email sudah digunakan",
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": fmt.Sprintf("Error database: %v", err),
		}
	}

	// Hash password
	hash, err := HashPassword(password)
	if err != nil {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": "Gagal mengenkripsi password",
		}
	}

	user := User{
		Username: username,
		Email:    email,
		Password: hash,
		Role:     "user", // default role
	}

	err = db.Create(&user).Error
	if err != nil {
		return map[string]interface{}{
			"Status":     "Gagal",
			"Keterangan": fmt.Sprintf("Gagal menyimpan user: %v", err),
		}
	}

	elapsed := time.Since(start).Seconds() * 1000
	return map[string]interface{}{
		"Status":           "Berhasil",
		"Keterangan":       "Registrasi berhasil",
		"Waktu_Registrasi": fmt.Sprintf("%.2f milisecond", elapsed),
	}
}

func Validasi(tokenString string, idUser uint) (bool, string) {
	// Parse token dan ambil claims
	claimsMap, err := ParseJWT(tokenString)
	if err != nil {
		return false, fmt.Sprintf("Token tidak valid: %v", err)
	}

	// Ambil user_id dari claims dan bandingkan dengan idUser
	// Klaim biasanya bertipe float64 saat di-unmarshal JSON
	claimUserID, ok := claimsMap["user_id"]
	if !ok {
		return false, "Token tidak mengandung user_id"
	}

	// Konversi ke uint (biasanya klaim berupa float64)
	var claimUserIDUint uint
	switch v := claimUserID.(type) {
	case float64:
		claimUserIDUint = uint(v)
	case int:
		claimUserIDUint = uint(v)
	case uint:
		claimUserIDUint = v
	default:
		return false, "Tipe user_id di token tidak dikenali"
	}

	if claimUserIDUint != idUser {
		return false, "User ID token tidak cocok dengan yang diberikan"
	}

	// Jika sampai sini, token valid dan user id cocok
	return true, "Validasi berhasil"
}
