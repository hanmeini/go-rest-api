package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go-flix-api/config"
	_ "go-flix-api/docs" // Import generated docs
	"go-flix-api/internal/auth"
	"go-flix-api/internal/middleware"
	"go-flix-api/internal/movie"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Flix API
// @version 1.0
// @description REST API untuk manajemen film dengan autentikasi JWT
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// ... (Komentar Swagger Anda) ...
func main() {
	// Setup logger, .env, config, dan database (tetap sama)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		slog.Warn("Peringatan: Gagal memuat file .env")
	}
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		slog.Error("Fatal: Gagal memuat config.yml", "error", err)
		os.Exit(1)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("Fatal: Gagal koneksi ke database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Koneksi database berhasil")
	slog.Info("Mencoba koneksi dengan DSN", "dsn", dsn)

	// === PERBAIKAN UTAMA DI SINI (Dependency Injection yang Benar) ===

	// 1. Inisialisasi semua service
	authService := auth.NewService(cfg)
	movieRepo := movie.NewRepository(db)
	movieService := movie.NewService(movieRepo)

	// 2. Inisialisasi semua handler, berikan service yang dibutuhkan
	authHandler := auth.NewHandler(authService)
	movieHandler := movie.NewHandler(movieService)

	// Router
	r := mux.NewRouter()

	// 3. Daftarkan rute dengan handler yang sudah diinisialisasi
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { /* ... */ }).Methods("GET")

	// Subrouter untuk Rute Terproteksi
	api := r.PathPrefix("/api").Subrouter()
	// 4. Berikan semua argumen yang dibutuhkan oleh middleware
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret, authService.IsTokenRevoked))

	api.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	api.HandleFunc("/movies", movieHandler.GetAllMovies).Methods("GET")
	api.HandleFunc("/movies", movieHandler.CreateMovie).Methods("POST", "OPTIONS")
	api.HandleFunc("/movies/{id}", movieHandler.GetMovieByID).Methods("GET")
	api.HandleFunc("/movies/{id}", movieHandler.UpdateMovie).Methods("PUT", "OPTIONS")
	api.HandleFunc("/movies/{id}", movieHandler.DeleteMovie).Methods("DELETE", "OPTIONS")

	// CORS Middleware dan Start Server (tetap sama)
	finalHandler := corsMiddleware(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	slog.Info("🚀 Server siap berjalan", "address", fmt.Sprintf("http://localhost:%s", port))
	if err := http.ListenAndServe(addr, finalHandler); err != nil {
		slog.Error("Gagal menjalankan server", "error", err)
		os.Exit(1)
	}
	slog.Info("Mencoba koneksi dengan DSN", "dsn", dsn)
}

// corsMiddleware (tetap sama)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
