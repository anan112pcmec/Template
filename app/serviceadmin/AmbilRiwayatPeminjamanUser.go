package serviceadmin

import (
	"fmt"

	"gorm.io/gorm"

)

func AmbilDataRiwayatPeminjamanUser(db *gorm.DB, nama, email string) []map[string]interface{} {
	var DataBukuDipinjam []Peminjamanbuku
	var hasilnya []map[string]interface{}

	if err := db.Unscoped().Table("peminjamanbuku").Where("namapeminjam = ?", nama).Find(&DataBukuDipinjam).Error; err != nil {
		fmt.Println("[ERROR] Gagal ambil data buku:", err)
	}

	if len(DataBukuDipinjam) == 0 {
		if err := db.Unscoped().Table("peminjamanbukus").Find(&DataBukuDipinjam).Error; err != nil {
			fmt.Println("[ERROR] Gagal ambil data buku dari fallback:", err)
			return nil
		}
	}

	fmt.Println("Ditemukan")
	count := 0

	for _, data := range DataBukuDipinjam {
		count++
		item := map[string]interface{}{
			"nomor":    count,
			"judul":    data.Judul,
			"tanggal":  data.Tanggal,
			"kategori": data.Kategori,
			"status":   data.Status,
		}
		hasilnya = append(hasilnya, item)
	}

	return hasilnya
}
