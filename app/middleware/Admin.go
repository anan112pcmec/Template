package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/serviceadmin"
)

// ============================
// Handler utama
// ============================

func AdminHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes := r.Context().Value(BodyKey)

		fmt.Println(bodyBytes)
		if bodyBytes == nil {
			http.Error(w, "Body tidak ditemukan di context", http.StatusBadRequest)
			return
		}

		bb, ok := bodyBytes.([]byte)
		if !ok {
			http.Error(w, "Tipe data body tidak valid", http.StatusInternalServerError)
			return
		}

		var data serviceadmin.RequestAdmin
		if err := json.Unmarshal(bb, &data); err != nil {
			http.Error(w, "Gagal parsing JSON RequestAdmin: "+err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("POST data diterima: %+v\n", data)

		var hasil any

		switch data.Tujuan {

		case "Memasukan Data Buku":
			var req serviceadmin.BukuBaruRequest
			if err := json.Unmarshal(bb, &req); err != nil {
				fmt.Println(err)
				http.Error(w, "Gagal parsing JSON BukuBaruRequest: "+err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println("ini req", req.Bahasa, req.Deskripsi, req.Harga)
			hasil = serviceadmin.MasukanBuku(db, &req)
		case "AmbilDataBuku":
			fmt.Println("AmbilDataBuku Jalan")
			hasil = serviceadmin.AmbilBukuAll(db)
		case "CommitEditBuku":
			var req serviceadmin.BukuBaruRequest
			if err := json.Unmarshal(bb, &req); err != nil {
				fmt.Println(err)
				http.Error(w, "Gagal parsing JSON BukuBaruRequest: "+err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println("CommitEditBukuDijalankan")
			hasil = serviceadmin.CommitEditBuku(db, &req)
		case "AmbilDataBukuRinci":
			fmt.Println("AmbilDataBukuRinci jalan")
			var req serviceadmin.BukuBaruRequest
			if err := json.Unmarshal(bb, &req); err != nil {
				fmt.Println(err)
				http.Error(w, "Gagal parsing JSON BukuBaruRequest: "+err.Error(), http.StatusBadRequest)
				return
			}
			hasil = serviceadmin.AmbilBukuRinci(db, req.ISBN, req.Jenis, req.Judul)

		case "Hapus Child Buku":
			fmt.Println("Hapus Child Buku Dijalankan")
			var req serviceadmin.BukuBaruRequest
			if gagal := json.Unmarshal(bb, &req); gagal != nil {
				fmt.Println(gagal, "terjadi kesalahan")
				http.Error(w, "Terjadi Kesalahan:"+gagal.Error(), http.StatusBadRequest)
				return
			}
			hasil = serviceadmin.HapusChildBuku(db, req.ISBN, req.ID)
		case "AmbilDataUserAdmin":
			fmt.Println("Mencoba Mengambil Data User Yang Ada")
			hasil = serviceadmin.AmbilDataUsers(db)

		case "AmbilDataUserRiwayatPeminjaman":
			fmt.Println("AmbilDataRiwayatPeminjamanUser")
			var req serviceadmin.UserRequest
			if gagal := json.Unmarshal(bb, &req); gagal != nil {
				fmt.Println(gagal, "terjadi kesalahan")
				http.Error(w, "Terjadi Kesalahan:"+gagal.Error(), http.StatusBadRequest)
				return
			}
			hasil = serviceadmin.AmbilDataRiwayatPeminjamanUser(db, req.Nama, req.Email)

		case "AmbilDataPeminjaman":
			fmt.Println("Menjalankan AmbilDataPeminjaman")
			var req serviceadmin.KontrolPeminjamanBuku
			if gagal := json.Unmarshal(bb, &req); gagal != nil {
				fmt.Println(gagal, "terjadi kesalahan")
				http.Error(w, "Terjadi Kesalahan:"+gagal.Error(), http.StatusBadRequest)
				return
			}
			hasil = serviceadmin.AmbilBukuDipinjam(db, req.Search)

		default:
			http.Error(w, "Tujuan tidak dikenali: "+data.Tujuan, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Data POST diproses di middleware",
			"tujuan":  data.Tujuan,
			"Hasil":   hasil,
		})
	}
}
