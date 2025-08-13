package serviceadmin

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilBukuRinci(db *gorm.DB, ISBN, jenis, judul string) []map[string]interface{} {
	var hasil []map[string]interface{}

	rows, err := db.Unscoped().
		Model(&models.BukuChild{}).
		Where("isbn = ? AND jenis = ? AND judul = ?", ISBN, jenis, judul).
		Rows()
	if err != nil {
		fmt.Println("Gagal ambil data BukuChild:", err)
		return hasil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Gagal ambil kolom:", err)
		return hasil
	}

	nomor := 1 // Nomor baris dimulai dari 1
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println("Gagal scan baris:", err)
			continue
		}

		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			rowMap[colName] = *val
		}

		// Tambahkan nomor unik per baris (bukan per kolom)
		rowMap["nomor"] = nomor
		nomor++

		hasil = append(hasil, rowMap)
	}

	// Kalau hasil kosong, kembalikan dummy data
	if len(hasil) == 0 {
		dummy := map[string]interface{}{
			"ID":         0,
			"Kode_induk": 0,
			"Judul":      "Dummy Judul",
			"Jenis":      "Dummy Jenis",
			"Harga":      int64(0),
			"Penulis":    "Dummy Penulis",
			"Penerbit":   "Dummy Penerbit",
			"Stok":       int64(0),
			"Tahun":      "0000",
			"ISBN":       "0000000000",
			"Kategori":   "Dummy Kategori",
			"Bahasa":     "Dummy Bahasa",
			"Status":     "Dummy Status",
			"Deskripsi":  "Dummy Deskripsi",
			"CreatedAt":  time.Now(),
			"UpdatedAt":  time.Now(),
			"DeletedAt":  nil,
			"nomor":      1,
		}
		hasil = append(hasil, dummy)
	}

	return hasil
}
