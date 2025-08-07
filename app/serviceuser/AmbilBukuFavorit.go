package serviceuser

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuFavoritDia(db *gorm.DB, nama, iduser string) []map[string]interface{} {
	var favoritnya []models.Favorit
	var hasil []map[string]interface{}

	idconv, err := strconv.Atoi(iduser)
	if err != nil {
		return []map[string]interface{}{
			{"kondisi": "gagal"},
		}
	}

	err = db.Where("nama_user = ? AND iduser = ?", nama, idconv).Find(&favoritnya).Error
	if err != nil {
		return []map[string]interface{}{
			{"kondisi": "gagal di find favorit"},
		}
	}

	for _, favorit := range favoritnya {
		var buku models.BukuInduk

		err := db.Unscoped().Where("judul = ? AND isbn = ?", favorit.JudulBuku, favorit.ISBN).First(&buku).Error
		if err == nil {

			var favoritCheck models.Favorit
			query := db.Where("iduser = ? AND isbn = ?", idconv, buku.ISBN).Find(&favoritCheck)

			var favoritkah string
			if query.RowsAffected > 0 {
				favoritkah = "buku disukai"
			} else {
				favoritkah = "buku tidak disukai"
			}

			// Siapkan data untuk append
			data := map[string]interface{}{
				"id":          buku.ID,
				"judul":       buku.Judul,
				"jenis":       buku.Jenis,
				"harga":       buku.Harga,
				"penulis":     buku.Penulis,
				"penerbit":    buku.Penerbit,
				"stok":        buku.Stok,
				"tahun":       buku.Tahun,
				"isbn":        buku.ISBN,
				"kategori":    buku.Kategori,
				"bahasa":      buku.Bahasa,
				"deskripsi":   buku.Deskripsi,
				"tujuan_aksi": buku.TujuanAksi,
				"diskon":      buku.Diskon,
				"rating":      buku.Rating,
				"disukai":     favoritkah,
			}

			if len(buku.Gambar) > 0 {
				data["gambar"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buku.Gambar)
			} else {
				data["gambar"] = nil
			}

			hasil = append(hasil, data)
		}
	}

	fmt.Println(hasil, "ini hasil dapet")
	return hasil
}
