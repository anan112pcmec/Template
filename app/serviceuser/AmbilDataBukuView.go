package serviceuser

import (
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilDataBukuView(db *gorm.DB, berdasar string) []map[string]interface{} {
	var bukuList []models.BukuInduk
	var result []map[string]interface{}

	// Ambil data berdasarkan filter
	switch berdasar {
	case "Popularitas":
		fmt.Println("Popularitas mau ngambil")
		db.Model(&models.BukuInduk{}).
			Order(`(rating * 0.6 + viewed / 1000 * 0.4) DESC`).
			Limit(8).
			Find(&bukuList)

	case "Rating Tertinggi":
		fmt.Println("Rating mau ngambil")
		db.Model(&models.BukuInduk{}).
			Order(`rating DESC`).
			Limit(8).
			Find(&bukuList)

	case "Jumlah Pembaca":
		fmt.Println("JumlahPembaca mau ngambil")
		db.Model(&models.BukuInduk{}).
			Order(`viewed DESC`).
			Limit(8).
			Find(&bukuList)
	}

	for _, buku := range bukuList {
		// Buat map manual hanya dengan field yang ingin ditampilkan
		temp := map[string]interface{}{
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
		}

		// Tambahkan gambar base64 jika ada
		if len(buku.Gambar) > 0 {
			temp["gambar"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buku.Gambar)
		} else {
			temp["gambar"] = nil
		}

		result = append(result, temp)
	}

	return result
}
