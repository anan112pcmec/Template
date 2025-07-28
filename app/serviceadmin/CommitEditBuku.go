package serviceadmin

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func CommitEditBuku(db *gorm.DB, data *BukuBaruRequest) string {
	fmt.Println("=== Mulai CommitEditBuku ===")
	fmt.Printf("Data diterima: %+v\n", data)

	// Cari buku dari tabel buku_induks berdasarkan ISBN
	var buku models.BukuInduk
	fmt.Println("Mencari buku dengan ISBN:", data.ISBN)
	if err := db.Unscoped().Table("buku_induks").Where("isbn = ?", data.ISBN).First(&buku).Error; err != nil {
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
		fmt.Println("📷 Mendeteksi gambar, mencoba decode...")
		parts := strings.Split(data.GambarBase64, ",")
		encoded := parts[len(parts)-1]
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			fmt.Println("⛔ Gagal decode gambar:", err)

			finalHarga, _ := strconv.Atoi(data.Harga)
			finalStok, _ := strconv.Atoi(data.Stok)

			// Siapkan field yang ingin diupdate
			updateFields := map[string]interface{}{
				"judul":       data.Judul,
				"jenis":       data.Jenis,
				"harga":       finalHarga,
				"penulis":     data.Penulis,
				"penerbit":    data.Penerbit,
				"stok":        finalStok,
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
		return "Gagal decode gambar: " + err.Error() + "Tapi Data Buku " + data.Judul + " Berhasil Diedit"
	}

	// Decode gambar jika ada

	finalHarga, _ := strconv.Atoi(data.Harga)
	finalStok, _ := strconv.Atoi(data.Stok)
	// Siapkan field yang ingin diupdate
	updateFields := map[string]interface{}{
		"judul":       data.Judul,
		"jenis":       data.Jenis,
		"harga":       finalHarga,
		"penulis":     data.Penulis,
		"penerbit":    data.Penerbit,
		"stok":        finalStok,
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

	updateFieldschild := map[string]interface{}{
		"judul":    data.Judul,
		"jenis":    data.Jenis,
		"harga":    finalHarga,
		"penulis":  data.Penulis,
		"penerbit": data.Penerbit,
		"stok":     finalStok,
		"tahun":    data.Tahun,
		"kategori": data.Kategori,
		"bahasa":   data.Bahasa,
	}

	if err := db.Unscoped().
		Table("buku_children").
		Where("isbn = ?", data.ISBN).
		Updates(updateFieldschild).Error; err != nil {
		fmt.Println("[ERROR] Gagal update data buku anak:", err)
	} else {
		fmt.Println("[INFO] Data buku anak berhasil diupdate.")
	}

	fmt.Println("✅ Update berhasil untuk buku:", data.Judul)
	fmt.Println("=== Selesai CommitEditBuku ===")

	return "Data Buku " + data.Judul + " Berhasil Diedit"
}
