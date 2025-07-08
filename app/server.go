package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anan112pcmec/Template/app/middleware"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type Appsetting struct {
	AppName, AppConf, AppPort string
}

type Dataconfig struct {
	dbHost, dbUser, dbPass, dbName, dbPort string
}

// enableCORS sama seperti sebelumnya
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", getenvi("ACCESS_CTRL", "Unauthorized"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// rateLimiter per IP menggunakan golang.org/x/time/rate

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients   = make(map[string]*clientLimiter)
	muClients sync.Mutex
)

// Mendapatkan limiter untuk IP tertentu
func getLimiter(ip string) *rate.Limiter {
	muClients.Lock()
	defer muClients.Unlock()

	cl, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 100) // rate 5 req/detik, burst 100
		clients[ip] = &clientLimiter{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// Update lastSeen
	cl.lastSeen = time.Now()
	return cl.limiter
}

// Cleanup clients yang sudah lama tidak aktif (misal > 5 menit)
func cleanupClients() {
	for {
		time.Sleep(time.Minute)

		muClients.Lock()
		for ip, cl := range clients {
			if time.Since(cl.lastSeen) > 120*time.Minute {
				delete(clients, ip)
			}
		}
		muClients.Unlock()
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP address", http.StatusInternalServerError)
			return
		}

		limiter := getLimiter(ip)
		if !limiter.Allow() {

			if r.Method == http.MethodPost {
				fmt.Println("🚨 Terjadi serangan ke POST:", r.URL.Path)
			} else if r.Method == http.MethodGet {
				fmt.Println("🚨 Terjadi serangan ke GET:", r.URL.Path)
			} else if r.URL.Path == "/ws" {
				fmt.Println("🚨 Terjadi serangan ke WebSocket")
			} else {
				fmt.Println("🚨 Terjadi serangan ke", r.Method, "pada", r.URL.Path)
			}

			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) initialize(appconfig Appsetting) {
	fmt.Println("Inisialisasi server:", appconfig.AppName)

	server.Router = mux.NewRouter()

	// Middleware CORS
	server.Router.Use(enableCORS)

	// Middleware rate limiter global
	server.Router.Use(rateLimitMiddleware)

	// Jalankan cleanup client limiter di background
	go cleanupClients()

	var dbConfig = Dataconfig{
		dbHost: getenvi("DBHOST", "localhost"),
		dbUser: getenvi("DBUSER", "postgres"),
		dbPass: getenvi("DBPASS", "Faiz"),
		dbName: getenvi("DBNAME", "perpustakaan"),
		dbPort: getenvi("DBPORT", "8082"),
	}

	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.dbHost, dbConfig.dbUser, dbConfig.dbPass, dbConfig.dbName, dbConfig.dbPort,
	)

	fmt.Println("DSN yang digunakan:", dsn)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Server Gagal Menyambungkan ke database")
	} else {
		fmt.Println("Berhasil terhubung ke database:", getenvi("DBNAME", "DBMU"))
		fmt.Println("DB_NAME dari .env:", os.Getenv("DB_NAME"))
	}

	var currentDB string
	server.DB.Raw("SELECT current_database();").Scan(&currentDB)
	fmt.Println("Database yang sedang digunakan:", currentDB)

	server.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo dari " + appconfig.AppName))
	})

	server.Router.Handle("/endpoint.go", middleware.PostHandler(server.DB)).Methods("POST", "OPTIONS")
	server.Router.Handle("/endpoint.go", middleware.GetHandler()).Methods("GET")
	server.Router.HandleFunc("/ws", middleware.HandleWebSocket)
}

func (server *Server) Run(alamat string) {
	fmt.Printf("Berjalan di port %s\n", alamat)
	log.Fatal(http.ListenAndServe(alamat, server.Router))
}

func getenvi(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error di env")
	}
	var server = Server{}
	var appconfig = Appsetting{
		AppName: getenvi("APPNAME", "backend"),
		AppConf: getenvi("APPENV", "developmentcoy"),
		AppPort: getenvi("APPPORT", "8081"),
	}

	server.initialize(appconfig)
	server.Run(appconfig.AppPort)
}
