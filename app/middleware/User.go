package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/serviceuser"
)

func UserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes := r.Context().Value(BodyKey)

		fmt.Println(bodyBytes)
		if bodyBytes == nil {
			http.Error(w, "Body Tidak Ada", http.StatusBadRequest)
			return
		}

		bb, ok := bodyBytes.([]byte)
		if !ok {
			http.Error(w, "Tipe data body tidak valid", http.StatusInternalServerError)
			return
		}

		var data serviceuser.RequestUser
		if err := json.Unmarshal(bb, &data); err != nil {
			http.Error(w, "Gagal parsing JSON RequestAdmin: "+err.Error(), http.StatusBadRequest)
			return
		}

		var hasil any

		switch data.Tujuan {
		case "AmbilDataBukuView":
			fmt.Println("AmbilDataBukuViewDijalankan")
			hasil = serviceuser.AmbilDataBukuView(db, data.Berdasarkan)
		default:
			http.Error(w, "Tujuan tidak dikenali: "+data.Tujuan, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Data POST diproses di middleware",
			"tujuan":    data.Tujuan,
			"HasilUser": hasil,
		})
	}
}
