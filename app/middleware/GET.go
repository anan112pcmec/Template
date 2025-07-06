package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handle Get Request

func GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetHandler dijalankan...")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Halo dari GET middleware handler",
		})
	}
}
