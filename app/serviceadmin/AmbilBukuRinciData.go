package serviceadmin

import (
	"fmt"

	"gorm.io/gorm"

)

func AmbilBukuRinci(db *gorm.DB, ISBN, jenis, judul string) []map[string]interface{} {
	var hasil []map[string]interface{}

	// Tambahkan .Unscoped() agar ambil semua, termasuk yang deleted_at IS NOT NULL
	rows, err := db.Unscoped().
		Model(&BukuChild{}).
		Where("isbn = ? AND jenis = ? AND judul = ?", ISBN, jenis, judul).
		Rows()
	if err != nil {
		fmt.Println("Gagal ambil data BukuChild:", err)
		return hasil
	}
	defer rows.Close()

	// Ambil nama kolom
	cols, _ := rows.Columns()
	for rows.Next() {
		// Buat penampung kolom
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan baris ke dalam kolom
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println("Gagal scan baris:", err)
			continue
		}

		// Masukkan ke map
		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			rowMap[colName] = *val
		}

		hasil = append(hasil, rowMap)
	}

	return hasil
}
