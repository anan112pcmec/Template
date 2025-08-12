package routines

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func UpDatabase(db *gorm.DB) {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Printf("Gagal migrasi Peminjaman: %v", err)
		} else {
			fmt.Println("Migrasi Peminjaman berhasil")
		}

		if err := db.AutoMigrate(&models.BukuInduk{}); err != nil {
			log.Printf("Gagal migrasi BukuInduk: %v", err)
		} else {
			fmt.Println("Migrasi BukuInduk berhasil")
		}

		if err := db.AutoMigrate(&models.BukuChild{}); err != nil {
			log.Printf("Gagal migrasi BukuChild: %v", err)
		} else {
			fmt.Println("Migrasi BukuChild berhasil")
		}

		if err := db.AutoMigrate(&models.PeminjamanBuku{}); err != nil {
			log.Printf("Gagal migrasi Peminjaman: %v", err)
		} else {
			fmt.Println("Migrasi Peminjaman berhasil")
		}

		if err := db.AutoMigrate(&models.Favorit{}); err != nil {
			log.Printf("Gagal migrasi Favorit: %v", err)
		} else {
			fmt.Println("Migrasi Favorit berhasil")
		}

		if err := db.AutoMigrate(&models.Komentar{}); err != nil {
			log.Printf("Gagal migrasi Komentar: %v", err)
		} else {
			fmt.Println("Migrasi Komentar berhasil")
		}
	}()

	wg.Wait()
	fmt.Println("Semua proses migrasi selesai.")
}

func CleanupTableBuku(db *gorm.DB) {

}
