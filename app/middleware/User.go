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

		fmt.Println(data.IdUser, "iniawalan")

		switch data.Tujuan {
		case "AmbilDataBukuView":
			fmt.Println("AmbilDataBukuViewDijalankan")
			hasil = serviceuser.AmbilDataBukuView(db, data.Berdasarkan)
		case "AmbilBukuScroll":
			fmt.Println("AmbilBukuScrollDijalankan")
			hasil = serviceuser.AmbilBukuScroll(db, data.BukanBuku, data.IdUser)
		case "AmbilDataBukuBerdasarkanPencarian":
			fmt.Println("AmbilDataBukuBerdasarkanPencarian dijalankan")
			hasil = serviceuser.AmbilDataBukuBerdasarkanPencarian(db, data.Pencarian, data.IdUser)
		case "AmbilBukuFavorit":
			fmt.Println("Pencegahan di case", data.IdUser, data.Favorit)
			hasil = serviceuser.AmbilBukuGenreFav(db, data.Favorit, data.IdUser)
		case "AmbilBukuFavoritdia":
			fmt.Println("AmbilBukuFavoritdia dijalankan")
			hasil = serviceuser.AmbilBukuFavoritDia(db, data.NamaUser, data.IdUser)
		case "Favorit":
			hasil = serviceuser.FavoritBuku(db, data.IdUser, data.ISBnBuku, data.NamaUser, data.NamaBuku)
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
