package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Template/app/serviceadmin"
)

// ============================
// Middleware untuk baca body
// ============================

type ctxKey string

const BodyKey ctxKey = "bodyBytes"

func BodyReaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Gagal membaca body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Reset kembali r.Body agar bisa dibaca ulang jika dibutuhkan
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Simpan body ke dalam context
		ctx := context.WithValue(r.Context(), BodyKey, bodyBytes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
