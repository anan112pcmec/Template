package serviceuser

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"

)

func FavoritBuku(db *gorm.DB, Id_user string, ISBN, Nama, Judul string) map[string]string {
	idnya, err := strconv.Atoi(Id_user)
	if err != nil {
		return map[string]string{
			"Status":  "Gagal Konversi Id: ID tidak valid",
			"Kondisi": "Tidak Disukai",
		}
	}

	var favorit models.Favorit
	err = db.Where("iduser = ? AND isbn = ?", idnya, ISBN).First(&favorit).Error

	if err == nil {
		if delErr := db.Unscoped().
			Where("iduser = ? AND isbn = ?", idnya, ISBN).
			Delete(&models.Favorit{}).Error; delErr != nil {
			return map[string]string{
				"Status":  "Gagal menghapus favorit",
				"Kondisi": "Masih Disukai",
			}
		}
		return map[string]string{
			"Status":  "Dihapus dari favorit",
			"Kondisi": "Tidak Disukai",
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newFavorit := models.Favorit{
			Id_User:   int64(idnya),
			ISBN:      ISBN,
			NamaUser:  Nama,
			JudulBuku: Judul,
		}
		if createErr := db.Create(&newFavorit).Error; createErr != nil {
			return map[string]string{
				"Status":  "Gagal menambahkan favorit",
				"Kondisi": "Tidak Disukai",
			}
		}
		return map[string]string{
			"Status":  "Berhasil ditambahkan ke favorit",
			"Kondisi": "Disukai",
		}
	} else {
		return map[string]string{
			"Status":  "Terjadi kesalahan saat pengecekan data",
			"Kondisi": "Tidak Diketahui",
		}
	}
}
