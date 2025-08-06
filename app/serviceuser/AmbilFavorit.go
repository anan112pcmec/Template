package serviceuser

import "gorm.io/gorm"

func AmbilFavorit(db *gorm.DB, Kondisi string) any {
	var kembalian any

	if Kondisi == "MengambilStatus" {

	} else if Kondisi == "MengambilBuku" {

	}

	return kembalian
}
