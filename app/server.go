package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	routines "github.com/anan112pcmec/Template/app/backend/Routines"
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

var (
	blockedIPs   = make(map[string]time.Time)
	muBlockedIPs sync.Mutex
)

type ipRequestInfo struct {
	count     int
	firstSeen time.Time
}

var (
	ipRequests   = make(map[string]*ipRequestInfo)
	muIpRequests sync.Mutex
)

func isBlocked(ip string) bool {
	muBlockedIPs.Lock()
	defer muBlockedIPs.Unlock()

	unblockTime, blocked := blockedIPs[ip]
	if !blocked {
		return false
	}

	if time.Now().After(unblockTime) {
		delete(blockedIPs, ip) // unblock jika waktu sudah lewat
		return false
	}
	return true
}

// Mendapatkan limiter untuk IP tertentu
func getLimiter(ip string) *rate.Limiter {
	muClients.Lock()
	defer muClients.Unlock()

	cl, exists := clients[ip]
	if !exists {
		reqLimitStr := getenvi("REQ_LIMIT", "5")
		burstLimitStr := getenvi("BURST_LIMIT", "100")

		reqLimitFloat, err := strconv.ParseFloat(reqLimitStr, 64)
		if err != nil {
			reqLimitFloat = 5
		}

		burstLimitInt, err := strconv.Atoi(burstLimitStr)
		if err != nil {
			burstLimitInt = 100
		}

		limiter := rate.NewLimiter(rate.Limit(reqLimitFloat), burstLimitInt)
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

		if isBlocked(ip) {
			fmt.Printf(" IP %s mencoba request saat masih diblokir\n", ip)
			return
		}

		muIpRequests.Lock()
		reqInfo, exists := ipRequests[ip]
		if !exists {
			reqInfo = &ipRequestInfo{count: 1, firstSeen: time.Now()}
			ipRequests[ip] = reqInfo
		} else {
			// Reset jika lebih dari 10 detik
			if time.Since(reqInfo.firstSeen) > 10*time.Second {
				reqInfo.count = 1
				reqInfo.firstSeen = time.Now()
			} else {
				reqInfo.count++
			}

			// Jika melebihi 150, blokir selama 5 menit
			if reqInfo.count > 150 {
				muBlockedIPs.Lock()
				blockedIPs[ip] = time.Now().Add(5 * time.Minute)
				muBlockedIPs.Unlock()

				delete(ipRequests, ip) // reset counter setelah blok
				muIpRequests.Unlock()

				fmt.Printf("🚫 IP %s diblokir karena mengirim %d request dalam 10 detik\n", ip, reqInfo.count)
				http.Error(w, "Terlalu banyak request. Anda diblokir selama 5 menit.", http.StatusTooManyRequests)
				return
			}
		}
		muIpRequests.Unlock()

		// Rate limiter biasa
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Rate limit terlampaui", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var jwtSecret = []byte("rahasia_kamu_ganti_in_production")

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Nama   string `json:"nama"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, nama string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Nama:   nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // token expired 2 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Fungsi cek apakah IP sedang diblokir

// Fungsi untuk blokir IP selama durasi tertentu
func blockIP(ip string, duration time.Duration) {
	muBlockedIPs.Lock()
	defer muBlockedIPs.Unlock()
	blockedIPs[ip] = time.Now().Add(duration)
	fmt.Printf("🚫 IP %s diblokir selama %v karena request mencurigakan\n", ip, duration)
}

func blockBadRequestsMiddleware(next http.Handler) http.Handler {
	// Regex pola umum serangan SQLi, XSS, Path Traversal, Command Injection
	patterns := []string{
		`(?i)\b(or|and)\b\s+\d+=\d+`, // or 1=1, and 2=2
		`(?i)union\s+select`,         // UNION SELECT
		`(?i)drop\s+table`,           // DROP TABLE
		`(?i)insert\s+into`,          // INSERT INTO
		`(?i)select\s+.+\s+from`,     // SELECT ... FROM
		`(?i)';--`,                   // Comment injection
		`(?i)sleep\(\d+\)`,           // sleep(10)

		// XSS
		"(?i)<script> .............</script>", // <script> ... </script>

		// Command Injection
		`(?i)(wget|curl|exec)\s`,
	}

	// Compile all regex
	var compiledPatterns []*regexp.Regexp
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		compiledPatterns = append(compiledPatterns, re)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback kalau SplitHostPort gagal
		}

		// Cek dulu apakah IP diblokir
		if isBlocked(ip) {
			// Silent drop atau kasih response
			http.Error(w, "Forbidden: IP diblokir sementara karena request mencurigakan", http.StatusForbidden)
			return
		}

		// Fungsi helper untuk cek string dengan semua pola
		checkPatterns := func(text string) bool {
			for _, re := range compiledPatterns {
				if re.MatchString(text) {
					return true
				}
			}
			return false
		}

		// Cek URL raw string
		if checkPatterns(r.URL.RawQuery) || checkPatterns(r.URL.Path) {
			blockIP(ip, 10*time.Minute)
			http.Error(w, "Forbidden: Request mengandung pola berbahaya", http.StatusForbidden)
			return
		}

		headersToCheck := []string{"User-Agent", "Referer", "Cookie"}
		for _, h := range headersToCheck {
			if checkPatterns(r.Header.Get(h)) {
				blockIP(ip, 10*time.Minute)
				http.Error(w, "Forbidden: Request mengandung pola berbahaya", http.StatusForbidden)
				return
			}
		}

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Body != nil {

				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					bodyStr := string(bodyBytes)
					if checkPatterns(bodyStr) {
						blockIP(ip, 10*time.Minute)
						http.Error(w, "Forbidden: Request mengandung pola berbahaya", http.StatusForbidden)
						return
					}

					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) initialize(appconfig Appsetting) {
	fmt.Println("Inisialisasi server:", appconfig.AppName)
	server.Router = mux.NewRouter()
	server.Router.Use(blockBadRequestsMiddleware)
	server.Router.Use(enableCORS)
	server.Router.Use(rateLimitMiddleware)

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

	go routines.UpDatabase(server.DB)

	server.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo dari " + appconfig.AppName))
	})

	server.Router.Handle("/admin", middleware.BodyReaderMiddleware(middleware.AdminHandler(server.DB)))
	server.Router.Handle("/user", middleware.BodyReaderMiddleware(middleware.UserHandler(server.DB)))
	server.Router.Handle("/auth", middleware.BodyReaderMiddleware(middleware.AuthHandle(server.DB)))
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
