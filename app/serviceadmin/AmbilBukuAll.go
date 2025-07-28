package serviceadmin

import (
	"encoding/base64"
	"fmt"

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
			"CreatedAt":  buku.CreatedAt,
			"UpdatedAt":  buku.UpdatedAt,
		}

		hasil = append(hasil, item)
	}

	return hasil
}
