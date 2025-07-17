package serviceadmin

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CommitEditBuku(db *gorm.DB, data *BukuBaruRequest) string {
	fmt.Println("=== Mulai CommitEditBuku ===")
	fmt.Printf("Data diterima: %+v\n", data)

	// Cari buku dari tabel buku_induks berdasarkan ISBN
	var buku BukuInduk
	fmt.Println("Mencari buku dengan ISBN:", data.ISBN)
	if err := db.Table("buku_induks").Where("isbn = ?", data.ISBN).First(&buku).Error; err != nil {
		fmt.Println("⛔ Buku tidak ditemukan:", err)
		return "Data Buku dengan ISBN " + data.ISBN + " tidak ditemukan"
	}
	fmt.Println("✅ Buku ditemukan:", buku.Judul)

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

	// Decode gambar jika ada
	if data.GambarBase64 != "" {
		fmt.Println("📷 Mendeteksi gambar, mencoba decode...")
		parts := strings.Split(data.GambarBase64, ",")
		encoded := parts[len(parts)-1]
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			fmt.Println("⛔ Gagal decode gambar:", err)

			// Siapkan field yang ingin diupdate
			updateFields := map[string]interface{}{
				"judul":       data.Judul,
				"jenis":       data.Jenis,
				"harga":       data.Harga,
				"penulis":     data.Penulis,
				"penerbit":    data.Penerbit,
				"stok":        data.Stok,
				"tahun":       data.Tahun,
				"kategori":    data.Kategori,
				"bahasa":      data.Bahasa,
				"deskripsi":   data.Deskripsi,
				"tujuan_aksi": data.Tujuan,
				"updated_at":  time.Now(),
			}
			fmt.Printf("🛠️ Field yang akan diupdate: %+v\n", updateFields)

			// Update di tabel buku_induks berdasarkan ISBN
			fmt.Println("🚀 Menjalankan update ke tabel buku_induks...")
			if err := db.Table("buku_induks").Where("isbn = ?", data.ISBN).Updates(updateFields).Error; err != nil {
				fmt.Println("⛔ Gagal update data:", err)
				return "Gagal mengedit buku: " + err.Error()
			}

		}
		gambarBytes = decoded
		fmt.Println("✅ Gambar berhasil didecode, panjang byte:", len(gambarBytes))
	} else {
		fmt.Println("ℹ️ Tidak ada gambar baru dikirim.")
	}

	// Siapkan field yang ingin diupdate
	updateFields := map[string]interface{}{
		"judul":       data.Judul,
		"jenis":       data.Jenis,
		"harga":       data.Harga,
		"penulis":     data.Penulis,
		"penerbit":    data.Penerbit,
		"stok":        data.Stok,
		"tahun":       data.Tahun,
		"kategori":    data.Kategori,
		"bahasa":      data.Bahasa,
		"deskripsi":   data.Deskripsi,
		"tujuan_aksi": data.Tujuan,
		"updated_at":  time.Now(),
	}

	if len(gambarBytes) > 0 {
		updateFields["gambar"] = gambarBytes
	}

	fmt.Printf("🛠️ Field yang akan diupdate: %+v\n", updateFields)

	// Update di tabel buku_induks berdasarkan ISBN
	fmt.Println("🚀 Menjalankan update ke tabel buku_induks...")
	if err := db.Table("buku_induks").Where("isbn = ?", data.ISBN).Updates(updateFields).Error; err != nil {
		fmt.Println("⛔ Gagal update data:", err)
		return "Gagal mengedit buku: " + err.Error()
	}

	fmt.Println("✅ Update berhasil untuk buku:", data.Judul)
	fmt.Println("=== Selesai CommitEditBuku ===")

	return "Data Buku " + data.Judul + " Berhasil Diedit"
}
