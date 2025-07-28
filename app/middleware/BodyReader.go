package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type ctxKey string

const BodyKey ctxKey = "bodyBytes"

func BodyReaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Gagal membaca body: "+err.Error(), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Simpan body ke dalam context
		ctx := context.WithValue(r.Context(), BodyKey, bodyBytes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
