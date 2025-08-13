package serviceuser

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func MasukanBukuKeDalamKeranjang(db *gorm.DB, ISBN, id_user string) map[string]string {
	fmt.Println("Mulai fungsi MasukanBukuKeDalamKeranjang")
	fmt.Println("ISBN diterima:", ISBN)
	fmt.Println("ID user diterima:", id_user)

	var ModelBuku models.BukuChild

	// konversi ID user dari string ke int
	id_user_final, err := strconv.Atoi(id_user)
	if err != nil {
		fmt.Println("Error konversi ID user:", err)
		return map[string]string{
			"status": "gagal",
			"error":  "ID user tidak valid",
		}
	}
	fmt.Println("ID user setelah konversi:", id_user_final)

	// cek buku dengan status "Ready" dan ISBN yang sesuai
	err = db.Unscoped().
		Where("status = ? AND isbn = ?", "Ready", ISBN).
		First(&ModelBuku).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Buku tidak ditemukan")
			return map[string]string{
				"status": "gagal",
				"error":  "Buku tidak tersedia",
			}
		}
		fmt.Println("Error query buku:", err)
		return map[string]string{
			"status": "gagal",
			"error":  "Terjadi kesalahan pada database",
		}
	}

	fmt.Println("Buku ditemukan:", ModelBuku)

	// buat struct Keranjang baru
	var keranjang = models.Keranjang{
		IDUser: int64(id_user_final), // konversi int -> int64
		IDBuku: int64(ModelBuku.ID),
		ISBN:   &ISBN,
	}

	fmt.Println("Struct keranjang yang akan dicek:", keranjang)

	// cek apakah buku sudah ada di keranjang
	var cekKeranjang models.Keranjang
	err = db.Where(`id_user = ? AND id_buku = ? AND "ISBN" = ?`, keranjang.IDUser, keranjang.IDBuku, keranjang.ISBN).
		First(&cekKeranjang).Error

	if err == nil {
		fmt.Println("Buku sudah ada di keranjang:", cekKeranjang)
		return map[string]string{
			"status": "gagal",
			"error":  "Buku sudah ada di keranjang",
		}
	} else if err != gorm.ErrRecordNotFound {
		fmt.Println("Error saat cek keranjang:", err)
		return map[string]string{
			"status": "gagal",
			"error":  "Terjadi kesalahan saat cek keranjang",
		}
	}

	fmt.Println("Buku belum ada di keranjang, akan ditambahkan")

	// masukkan ke database
	if err := db.Create(&keranjang).Error; err != nil {
		fmt.Println("Error saat menambahkan ke keranjang:", err)
		return map[string]string{
			"status": "gagal",
			"error":  "Gagal menambahkan buku ke keranjang",
		}
	}

	fmt.Println("Berhasil menambahkan buku ke keranjang:", keranjang)

	return map[string]string{
		"status": "berhasil",
	}
}
