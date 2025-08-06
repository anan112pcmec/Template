package serviceuser

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"

)

func AmbilBukuGenreFav(db *gorm.DB, favorit []string, iduser string) []map[string]interface{} {
	var hasil []models.BukuInduk
	var output []map[string]interface{}
	fmt.Println(iduser, "takdeee")

	idbeneran, gagal := strconv.Atoi(iduser)
	if gagal != nil {
		return []map[string]interface{}{
			{
				"status": "gagalcuk",
			},
		}
	}

	err := db.Where("kategori IN ?", favorit).Find(&hasil).Error
	if err != nil {
		return nil
	}

	for _, buku := range hasil {
		item := map[string]interface{}{
			"id":          buku.ID,
			"judul":       buku.Judul,
			"jenis":       buku.Jenis,
			"harga":       buku.Harga,
			"penulis":     buku.Penulis,
			"penerbit":    buku.Penerbit,
			"stok":        buku.Stok,
			"tahun":       buku.Tahun,
			"isbn":        buku.ISBN,
			"kategori":    buku.Kategori,
			"bahasa":      buku.Bahasa,
			"deskripsi":   buku.Deskripsi,
			"tujuan_aksi": buku.TujuanAksi,
			"diskon":      buku.Diskon,
			"rating":      buku.Rating,
		}

		var favorit models.Favorit
		query := db.Where("iduser = ? AND isbn = ?", idbeneran, buku.ISBN).Find(&favorit)

		var favoritkah string
		if query.RowsAffected > 0 {
			favoritkah = "buku disukai"
		} else {
			favoritkah = "buku tidak disukai"
		}

		item["disukai"] = favoritkah

		if len(buku.Gambar) > 0 {
			item["gambar"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buku.Gambar)
		} else {
			item["gambar"] = nil
		}

		output = append(output, item)
	}

	return output
}
