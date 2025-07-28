package serviceadmin

import (
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuDipinjam(db *gorm.DB, Dicari string) []map[string]interface{} {
	var bukuList []models.BukuChild
	var hasil []map[string]interface{}

	// Ambil semua buku dengan status 'dipinjam'
	if Dicari == "Masa_peminjaman" {
		if err := db.Unscoped().
			Table("buku_children").
			Where("status = ?", "dipinjam").
			Find(&bukuList).Error; err != nil {
			fmt.Println("[ERROR] Gagal ambil data buku:", err)
			return nil
		}
	}

	if Dicari == "Semua" {
		if err := db.Unscoped().Table("buku_children").Where("status != ?", "Ready").Find(&bukuList).Error; err != nil {
			fmt.Println("[ERROR] Gagal ambil data buku:", err)
			return nil
		}
	}

	if Dicari == "Belum_Dikembalikan" {
		if err := db.Unscoped().Table("buku_children").Where("status = ?", "Belum Dikembalikan").Find(&bukuList).Error; err != nil {
			fmt.Println("[ERROR] Gagal Cyukk", err)
			return nil
		}
	}

	if Dicari == "Dikembalikan_Hari_Ini" {
		if err := db.Unscoped().Table("buku_children").Where("status = ?", "Dikembalikan Hari Ini").Find(&bukuList).Error; err != nil {
			fmt.Println("[ERROR] Gagal Cyukk", err)
			return nil
		}
	}

	if Dicari == "Dikembalikan" {
		if err := db.Unscoped().Table("buku_children").Where("status = ?", "Dikembalikan").Find(&bukuList).Error; err != nil {
			fmt.Println("[ERROR] Gagal Cyukk", err)
			return nil
		}
	}

	count := 0
	for _, buku := range bukuList {
		var gambar string
		count++
		// Ambil gambar dari tabel buku_induks berdasarkan ISBN dan Judul
		if err := db.Unscoped().
			Table("buku_induks").
			Select("gambar").
			Where("isbn = ? AND judul = ?", buku.ISBN, buku.Judul).
			Scan(&gambar).Error; err != nil {
			fmt.Println("[ERROR] Gagal ambil gambar buku:", err)
			return nil
		}

		var gambarBase64 string
		if gambar != "" {
			gambarBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(gambar))
		}

		item := map[string]interface{}{
			"nomor":     count,
			"Status":    buku.Status,
			"Judul":     buku.Judul,
			"Kode":      buku.ID,
			"Jenis":     buku.Jenis,
			"Harga":     buku.Harga,
			"Penulis":   buku.Penulis,
			"Penerbit":  buku.Penerbit,
			"Stok":      buku.Stok,
			"Tahun":     buku.Tahun,
			"ISBN":      buku.ISBN,
			"Kategori":  buku.Kategori,
			"Bahasa":    buku.Bahasa,
			"Deskripsi": buku.Deskripsi,
			"Gambar":    gambarBase64,
			"CreatedAt": buku.CreatedAt,
			"UpdatedAt": buku.UpdatedAt,
		}

		hasil = append(hasil, item)
	}

	return hasil
}
