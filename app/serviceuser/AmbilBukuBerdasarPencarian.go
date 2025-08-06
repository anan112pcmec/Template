package serviceuser

import (
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilDataBukuBerdasarkanPencarian(db *gorm.DB, Bukudicari, iduser string) []map[string]interface{} {
	var data []models.BukuInduk

	fmt.Println("Nyoba nyari judul dulu buat pencarian")

	if Bukudicari == "Menunjukan Semua Buku" {
		err := db.Unscoped().Select("*").Find(&data).Error
		if err != nil {
			return nil
		}
	} else {
		err := db.Unscoped().Select("*").Where("judul ILIKE ?", "%"+Bukudicari+"%").Find(&data).Error
		if err != nil {
			return nil
		}
	}

	fmt.Println("Judulnya ketemu lanjut parsing")

	var hasil []map[string]interface{}
	for _, buku := range data {
		bukuMap := map[string]interface{}{
			"id":        buku.ID,
			"judul":     buku.Judul,
			"jenis":     buku.Jenis,
			"harga":     buku.Harga,
			"penulis":   buku.Penulis,
			"penerbit":  buku.Penerbit,
			"stok":      buku.Stok,
			"tahun":     buku.Tahun,
			"isbn":      buku.ISBN,
			"kategori":  buku.Kategori,
			"bahasa":    buku.Bahasa,
			"deskripsi": buku.Deskripsi,
			"rating":    buku.Rating,
			"viewed":    buku.Viewed,
			"diskon":    buku.Diskon,
		}

		var favorit models.Favorit
		query := db.Where("iduser = ? AND isbn = ?", iduser, buku.ISBN).Find(&favorit)

		var favoritkah string
		if query.RowsAffected > 0 {
			favoritkah = "buku disukai"
		} else {
			favoritkah = "buku tidak disukai"
		}

		bukuMap["disukai"] = favoritkah

		if len(buku.Gambar) > 0 {
			bukuMap["gambar"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buku.Gambar)
		} else {
			bukuMap["gambar"] = nil
		}

		hasil = append(hasil, bukuMap)
	}

	return hasil
}
