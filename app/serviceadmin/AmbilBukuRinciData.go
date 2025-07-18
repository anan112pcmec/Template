package serviceadmin

import (
	"fmt"

	"gorm.io/gorm"
)

func AmbilBukuRinci(db *gorm.DB, ISBN, jenis, judul string) []map[string]interface{} {
	var hasil []map[string]interface{}

	rows, err := db.Unscoped().
		Model(&BukuChild{}).
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

	return hasil
}
