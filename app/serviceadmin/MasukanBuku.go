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
)

// Fungsi untuk input dan migrasi otomatis

func randomSixDigit() int {
	rand.Seed(time.Now().UnixNano())  // Seed random dengan waktu sekarang
	return rand.Intn(900000) + 100000 // 0..899999 + 100000 = 100000..999999
}

func MasukanBuku(db *gorm.DB, data *BukuBaruRequest) string {
	fmt.Println("Ini dari service Masukan Buku attempt ke-", data.Bahasa, data.Harga, data.Judul)

	// Auto migrate tabel
	if err := db.AutoMigrate(&BukuInduk{}); err != nil {
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

	buku := BukuInduk{
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
	err1 := db.Select("isbn").Where("isbn = ?", buku.ISBN).First(&isbnConfirm).Error

	if err1 != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Ada kesalahan pada verifikasi data, akan retry.")
		// Tambahkan retry logic di sini jika diperlukan
	}

	if isbnConfirm != 0 {
		return "Buku dengan ISBN ini sudah ada di DB. Input buku lain."
	}

	// Simpan ke DB
	if err := db.Create(&buku).Error; err != nil {
		return "Gagal menyimpan Buku ISBN Buku Ini Telah Digunakan Buku Lain Dan tak mungkin Valid"
	}

	var idConfirm int64
	err2 := db.Select("id").Where("ISBN = ?", buku.ISBN).First(idConfirm).Error
	if err2 != nil {
		fmt.Println("Ada kesalahan pada verifikasi data, akan retry. Attempt ke:")

		// Hapus data berdasarkan isbn yang dikirim
		if delErr := db.Where("isbn = ?", data.ISBN).Delete(&BukuInduk{}).Error; delErr != nil {
			fmt.Println("Gagal hapus data saat retry:", delErr.Error())
		}
	}
	for i := 0; i < Stokfinal; i++ {
		bukuChild := BukuChild{
			ID:        uint(randomSixDigit()),
			KodeInduk: uint(idConfirm),
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

		if err := db.Create(&bukuChild).Error; err != nil {
			return "Gagal menyimpan buku: " + err.Error()
		}

	}

	return "Data Buku " + data.Judul + " Berhasil Dimasukan"
}
