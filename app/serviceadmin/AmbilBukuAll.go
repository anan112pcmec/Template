package serviceadmin

import (
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuAll(db *gorm.DB) []map[string]interface{} {
	var bukuList []models.BukuInduk
	var hasil []map[string]interface{}

	// Ambil semua data dari tabel buku_induks
	if err := db.Unscoped().Table("buku_induks").Find(&bukuList).Error; err != nil {
		fmt.Println("[ERROR] Gagal ambil data buku:", err)
		return nil
	}

	for _, buku := range bukuList {
		var gambarBase64 string
		if len(buku.Gambar) > 0 {
			// Tambahkan prefix MIME type PNG
			gambarBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buku.Gambar)
		}

		item := map[string]interface{}{
			"Judul":      buku.Judul,
			"Jenis":      buku.Jenis,
			"Harga":      buku.Harga,
			"Penulis":    buku.Penulis,
			"Penerbit":   buku.Penerbit,
			"Stok":       buku.Stok,
			"Tahun":      buku.Tahun,
			"ISBN":       buku.ISBN,
			"Kategori":   buku.Kategori,
			"Bahasa":     buku.Bahasa,
			"Deskripsi":  buku.Deskripsi,
			"TujuanAksi": buku.TujuanAksi,
			"Gambar":     gambarBase64,
			"Diskon":     buku.Diskon,
			"CreatedAt":  buku.CreatedAt,
			"UpdatedAt":  buku.UpdatedAt,
			"Rating":     buku.Rating,
			"Viewed":     buku.Viewed,
		}

		hasil = append(hasil, item)
	}

	// Kalau kosong, buat 2 dummy data
	if len(hasil) == 0 {
		for i := 1; i <= 2; i++ {
			dummy := map[string]interface{}{
				"Judul":      fmt.Sprintf("Dummy Judul %d", i),
				"Jenis":      "Dummy Jenis",
				"Harga":      int64(0),
				"Penulis":    "Dummy Penulis",
				"Penerbit":   "Dummy Penerbit",
				"Stok":       int64(0),
				"Tahun":      "0000",
				"ISBN":       fmt.Sprintf("000000000%d", i),
				"Kategori":   "Dummy Kategori",
				"Bahasa":     "Dummy Bahasa",
				"Deskripsi":  "Dummy Deskripsi",
				"TujuanAksi": "Dummy Tujuan",
				"Gambar":     "",
				"Diskon":     float64(0),
				"CreatedAt":  time.Now(),
				"UpdatedAt":  time.Now(),
				"Rating":     float64(0),
				"Viewed":     int64(0),
			}
			hasil = append(hasil, dummy)
		}
	}

	return hasil
}
