package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// Handle Post Request

func PostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("PostHandler dijalankan...")

		var data RequestData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Gagal parsing JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Contoh logika (kamu bisa sesuaikan dengan fungsimu seperti user.Login dll)
		fmt.Printf("POST data diterima: %+v\n", data)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Data POST diproses di middleware",
			"tujuan":  data.Tujuan,
		})
	}
}
