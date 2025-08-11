package serviceuser

import (
	"encoding/base64"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuBestSeller(db *gorm.DB) []map[string]interface{} {
	kategoriList := []string{"Fiksi", "Non-Fiksi", "Sejarah", "Lainnya"}
	hasil := []map[string]interface{}{}

	for _, kategori := range kategoriList {
		var bukuList []models.BukuInduk

		err := db.Unscoped().
			Where("kategori = ?", kategori).
			Order("penjualan DESC").
			Limit(8).
			Find(&bukuList).Error

		if err != nil {
			hasil = append(hasil, map[string]interface{}{
				"Jenis": kategori,
				"Buku": map[string]interface{}{
					"judul":  "",
					"status": "Gagal",
					"pesan":  err.Error(),
				},
			})
			continue
		}

		for _, b := range bukuList {
			var gambarBase64 string
			if len(b.Gambar) > 0 {
				gambarBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Gambar)
			}

			hasil = append(hasil, map[string]interface{}{
				"Jenis": kategori,
				"Buku": map[string]interface{}{
					"id":          b.ID,
					"judul":       b.Judul,
					"jenis":       b.Jenis,
					"harga":       b.Harga,
					"penulis":     b.Penulis,
					"penerbit":    b.Penerbit,
					"stok":        b.Stok,
					"tahun":       b.Tahun,
					"isbn":        b.ISBN,
					"kategori":    b.Kategori,
					"bahasa":      b.Bahasa,
					"deskripsi":   b.Deskripsi,
					"tujuan_aksi": b.TujuanAksi,
					"gambar":      gambarBase64,
					"diskon":      b.Diskon,
					"rating":      b.Rating,
					"viewed":      b.Viewed,
					"disukai":     "",
				},
			})
		}
	}

	return hasil
}
