package routines

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/backend/models"
)

func DownDatabase(db *gorm.DB) {
	var wg sync.WaitGroup

	wg.Add(3)

	// Drop User (bukan relasi foreign, jadi bisa duluan atau belakangan)
	go func() {
		defer wg.Done()
		if err := db.Migrator().DropTable(&models.User{}); err != nil {
			log.Printf("Gagal drop tabel User: %v", err)
		} else {
			fmt.Println("Tabel User berhasil di-drop")
		}
	}()

	// Drop BukuChild dulu karena biasanya punya foreign key ke BukuInduk
	go func() {
		defer wg.Done()
		if err := db.Migrator().DropTable(&models.BukuChild{}); err != nil {
			log.Printf("Gagal drop tabel BukuChild: %v", err)
		} else {
			fmt.Println("Tabel BukuChild berhasil di-drop")
		}
	}()

	// Drop BukuInduk setelah BukuChild
	go func() {
		defer wg.Done()
		if err := db.Migrator().DropTable(&models.BukuInduk{}); err != nil {
			log.Printf("Gagal drop tabel BukuInduk: %v", err)
		} else {
			fmt.Println("Tabel BukuInduk berhasil di-drop")
		}
	}()

	wg.Wait()
	fmt.Println("Semua tabel berhasil di-drop (DownDatabase selesai).")
}
