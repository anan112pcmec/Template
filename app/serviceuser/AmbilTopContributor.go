package serviceuser

import (
	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func AmbilTopContributor(db *gorm.DB, jenis string) []map[string]interface{} {
	var komentar []models.Komentar

	err := db.
		Unscoped().
		Where("kepada_jenis = ?", jenis).
		Order("id ASC").
		Limit(100).
		Find(&komentar).Error

	if err != nil {
		return []map[string]interface{}{
			{
				"Status": "GagalCuk",
			},
		}
	}

	// Konversi slice komentar ke slice map[string]interface{}
	result := make([]map[string]interface{}, 0, len(komentar))

	for _, k := range komentar {
		m := map[string]interface{}{
			"id":           k.ID,
			"id_user":      k.IdUser,
			"komentar":     k.Komentar,
			"namabuku":     k.NamaBuku,
			"isbn":         k.ISBN,
			"namauser":     k.NamaUser,
			"bintang":      k.Bintang,
			"kepada_jenis": k.KepadaJenis,
			"dalam_jenis":  k.DalamJenis,
		}
		result = append(result, m)
	}

	return result
}
