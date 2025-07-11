package userauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret digunakan untuk signing
var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecret = []byte(secret)
		log.Println("JWT secret loaded from environment variable")
	} else {
		key := make([]byte, 32) // HS256 butuh 256-bit = 32 bytes
		_, err := rand.Read(key)
		if err != nil {
			log.Fatalf("Failed to generate JWT secret: %v", err)
		}
		jwtSecret = key
		log.Printf("JWT secret generated automatically (base64): %s\n", base64.StdEncoding.EncodeToString(jwtSecret))
	}
}

// CreateJWT membuat token dari klaim key-value (map)
func CreateJWT(claimsMap map[string]interface{}, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	for key, value := range claimsMap {
		claims[key] = value
	}
	// Set expiry dan issued at
	now := time.Now()
	claims["exp"] = now.Add(expiresIn).Unix()
	claims["iat"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT membaca token dan mengembalikan claims map
func ParseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak dikenali: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid")
}

// IsJWTValid hanya cek validitas token saja
func IsJWTValid(tokenString string) bool {
	_, err := ParseJWT(tokenString)
	return err == nil
}
