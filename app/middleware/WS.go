package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler WebSocket yang membaca JSON dan parsing ke struct
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket gagal upgrade:", err)
		http.Error(w, "Tidak dapat membuka koneksi websocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("Koneksi WebSocket terbuka")

	for {
		// Terima pesan dari client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Koneksi WebSocket terputus:", err)
			break
		}

		// Unmarshal ke struct RequestData
		var reqData RequestData
		err = json.Unmarshal(message, &reqData)
		if err != nil {
			log.Println("Gagal parsing JSON:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Format data salah"))
			continue
		}

		log.Printf("📩 Data diterima: Tujuan=%v, Nama=%s, Pw=%s", reqData.Tujuan, reqData.Nama, reqData.Pw)

		// Kirim respons balik (opsional)
		response := map[string]string{
			"status": "ok",
			"pesan":  "Data diterima",
			"nama":   reqData.Nama,
		}
		respJSON, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, respJSON)
	}
}
