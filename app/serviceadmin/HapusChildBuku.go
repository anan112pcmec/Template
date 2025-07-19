package serviceadmin

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

)

func HapusChildBuku(db *gorm.DB, ISBN, kodeId string) string {
	IdFinal, err := strconv.Atoi(kodeId)
	if err != nil {
		return "❌ Kode ID tidak valid: " + err.Error()
	}

	// Ambil judul dari buku_induks berdasarkan ISBN
	var judul string
	if err := db.Table("buku_induks").
		Select("judul").
		Where("isbn = ?", ISBN).
		Scan(&judul).Error; err != nil {
		fmt.Println("⛔ Gagal mengambil judul:", err)
		return "Gagal menghapus: " + err.Error()
	}

	if judul == "" {
		return "❌ Data buku induk tidak ditemukan."
	}

	// Hapus child buku berdasarkan ISBN, judul, dan ID
	if err := db.Table("buku_children").
		Where("isbn = ? AND judul = ? AND id = ?", ISBN, judul, IdFinal).
		Delete(nil).Error; err != nil {
		fmt.Println("⛔ Gagal menghapus data child:", err)
		return "Gagal menghapus: " + err.Error()
	}

	return "Data buku " + judul + " Dengan Kode " + kodeId + " berhasil dihapus dari sistem."

}
