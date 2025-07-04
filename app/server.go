package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anan112pcmec/gomodultem/app/middleware"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type Appsetting struct {
	AppName, AppConf, AppPort string
}

type Dataconfig struct {
	dbHost, dbUser, dbPass, dbPort string
}

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

func (server *Server) initialize(appconfig Appsetting) {
	fmt.Println("Inisialisasi server:", appconfig.AppName)
	server.Router = mux.NewRouter()
	server.Router.Use(enableCORS)

	var dbConfig = Dataconfig{
		dbHost: getenvi("DBHOST", "localhost"),
		dbUser: getenvi("DBUSER", "postgres"),
		dbPass: getenvi("DBPASS", "Faiz"),
		dbPort: getenvi("DBPORT", "8082"),
	}

	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=perpustakaanfaiz port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.dbHost, dbConfig.dbUser, dbConfig.dbPass, dbConfig.dbPort,
	)

	fmt.Println("DSN yang digunakan:", dsn)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Server Gagal Menyambungkan ke database")
	} else {
		fmt.Println("Berhasil terhubung ke database:", "kasir_go")
		fmt.Println("DB_NAME dari .env:", os.Getenv("DB_NAME"))
	}

	var currentDB string
	server.DB.Raw("SELECT current_database();").Scan(&currentDB)
	fmt.Println("Database yang sedang digunakan:", currentDB)

	// Root route
	server.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo dari " + appconfig.AppName))
	})

	// ⏬ Ganti handler manual dengan handler dari package middleware ⏬
	server.Router.Handle("/endpoint.go", middleware.PostHandler(server.DB)).Methods("POST", "OPTIONS")
	server.Router.Handle("/endpoint.go", middleware.GetHandler()).Methods("GET")
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
