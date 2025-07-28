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

	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := db.AutoMigrate(&models.BukuInduk{}); err != nil {
			log.Printf("Gagal migrasi BukuInduk: %v", err)
		} else {
			fmt.Println("Migrasi BukuInduk berhasil")
		}
	}()

	// Migrasi BukuChild
	go func() {
		defer wg.Done()
		if err := db.AutoMigrate(&models.BukuChild{}); err != nil {
			log.Printf("Gagal migrasi BukuChild: %v", err)
		} else {
			fmt.Println("Migrasi BukuChild berhasil")
		}
	}()

	go func() {
		defer wg.Done()
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Printf("Gagal migrasi User: %v", err)
		} else {
			fmt.Println("Migrasi User berhasil")
		}
	}()

	go func() {
		defer wg.Done()
		if err := db.AutoMigrate(&models.PeminjamanBuku{}); err != nil {
			log.Printf("Gagal migrasi Peminjaman: %v", err)
		} else {
			fmt.Println("Migrasi Peminjamanberhasil")
		}
	}()

	wg.Wait()
	fmt.Println("Semua proses migrasi selesai.")
}

func CleanupTableBuku(db *gorm.DB) {

}
