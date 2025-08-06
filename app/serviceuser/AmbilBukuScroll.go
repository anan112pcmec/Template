package serviceuser

import (
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuScroll(db *gorm.DB, jangandiambil []string, iduser string) []map[string]interface{} {
	var datanya []models.BukuInduk
	var hasil []map[string]interface{}

	fmt.Println("Ambilbukuscrolljalan")

	// Ambil lebih banyak data untuk memberi ruang penyaringan manual
	err := db.Unscoped().Order("RANDOM()").Limit(20).Find(&datanya).Error
	if err != nil {
		fmt.Println("Gagal mengambil data:", err)
		return nil
	}

	judulTolak := make(map[string]bool)
	for _, j := range jangandiambil {
		judulTolak[j] = true
	}

	// Filter secara manual
	counter := 0
	for _, buku := range datanya {
		if judulTolak[buku.Judul] {
			continue
		}

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
		counter++

		if counter == 8 {
			break
		}
	}

	return hasil
}
