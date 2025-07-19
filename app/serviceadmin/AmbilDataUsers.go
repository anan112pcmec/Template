package serviceadmin

import (
	"fmt"

	"gorm.io/gorm"
)

func AmbilDataUsers(db *gorm.DB) []map[string]interface{} {

	var Users []User
	var hasilnya []map[string]interface{}

	if err := db.Unscoped().Table("users").Find(&Users).Error; err != nil {
		fmt.Println("[ERROR] Gagal ambil data buku:", err)
		return nil
	}

	fmt.Println(Users)
	count := 0

	for _, user := range Users {
		count++
		item := map[string]interface{}{
			"nomor":      count,
			"id":         user.ID,
			"nama":       user.Nama,
			"favorit":    user.Favorit,
			"kreditskor": user.KreditSkor,
			"email":      user.Email,
			"alamat":     user.Alamat,
			"status":     user.Status,
			"bergabung":  user.Bergabung,
		}
		hasilnya = append(hasilnya, item)
	}

	return hasilnya
}
