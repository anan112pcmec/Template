package serviceadmin

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

// Fungsi untuk input dan migrasi otomatis

func randomSixDigit() int {
	rand.Seed(time.Now().UnixNano())  // Seed random dengan waktu sekarang
	return rand.Intn(900000) + 100000 // 0..899999 + 100000 = 100000..999999
}

func MasukanBuku(db *gorm.DB, data *BukuBaruRequest) string {
	fmt.Println("Ini dari service Masukan Buku attempt ke-", data.Bahasa, data.Harga, data.Judul)

	// Auto migrate tabel
	if err := db.AutoMigrate(&models.BukuInduk{}); err != nil {
		return "Gagal auto-migrasi tabel BukuInduk: " + err.Error()
	}

	// Bersihkan prefix base64
	cleanedBase64 := data.GambarBase64
	if strings.HasPrefix(cleanedBase64, "data:") {
		parts := strings.SplitN(cleanedBase64, ",", 2)
		if len(parts) == 2 {
			cleanedBase64 = parts[1]
		}
	}

	// Decode base64 menjadi raw byte
	gambarBytes, err := base64.StdEncoding.DecodeString(cleanedBase64)
	if err != nil {
		return "Gagal decode gambar: " + err.Error()
	}

	hargafinal, _ := strconv.Atoi(data.Harga)
	Stokfinal, _ := strconv.Atoi(data.Stok)

	IdInduk := randomSixDigit()

	buku := models.BukuInduk{
		ID:        uint(IdInduk),
		Judul:     data.Judul,
		Jenis:     data.Jenis,
		Harga:     int64(hargafinal),
		Penulis:   data.Penulis,
		Penerbit:  data.Penerbit,
		Stok:      int64(Stokfinal),
		Tahun:     data.Tahun,
		ISBN:      data.ISBN,
		Kategori:  data.Kategori,
		Bahasa:    data.Bahasa,
		Deskripsi: data.Deskripsi,
		Gambar:    gambarBytes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var isbnConfirm int64 = 0

	Transaction := db.Transaction(func(tx *gorm.DB) error {
		err1 := tx.Model(&models.BukuInduk{}).Select("isbn").Where("isbn = ?", buku.ISBN).Scan(&isbnConfirm).Error
		if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
			return fmt.Errorf("gagal verifikasi ISBN: %w", err1)
		}
		if isbnConfirm != 0 {
			return fmt.Errorf("buku dengan ISBN ini sudah ada di DB")
		}

		// Simpan BukuInduk
		if err := tx.Create(&buku).Error; err != nil {
			return fmt.Errorf("gagal menyimpan BukuInduk: %w", err)
		}

		// Ambil ID induk untuk BukuChild
		var idConfirm uint
		err2 := tx.Model(&models.BukuInduk{}).Select("id").Where("isbn = ?", buku.ISBN).Scan(&idConfirm).Error
		if err2 != nil {
			// Rollback eksplisit dengan menghapus yang sudah dimasukkan (opsional karena transaksi akan rollback juga)
			_ = tx.Where("isbn = ?", buku.ISBN).Delete(&models.BukuInduk{})
			return fmt.Errorf("gagal verifikasi ID induk: %w", err2)
		}

		// Buat BukuChild sebanyak stok
		for i := 0; i < Stokfinal; i++ {
			bukuChild := models.BukuChild{
				ID:        uint(randomSixDigit()),
				KodeInduk: idConfirm,
				Judul:     data.Judul,
				Jenis:     data.Jenis,
				Harga:     int64(hargafinal),
				Penulis:   data.Penulis,
				Penerbit:  data.Penerbit,
				Stok:      int64(Stokfinal),
				Tahun:     data.Tahun,
				ISBN:      data.ISBN,
				Kategori:  data.Kategori,
				Bahasa:    data.Bahasa,
				Status:    "Ready",
				Deskripsi: data.Deskripsi,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := tx.Create(&bukuChild).Error; err != nil {
				return fmt.Errorf("gagal menyimpan BukuChild ke-%d: %w", i+1, err)
			}
		}
		return nil // sukses, akan di-commit
	})

	if Transaction != nil {
		return "Terjadi kesalahan: " + err.Error()
	}

	return "Data Buku " + data.Judul + " Berhasil Dimasukan"
}
