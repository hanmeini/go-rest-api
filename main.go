package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"go-flix-api/auth"
	"go-flix-api/handlers"
	"go-flix-api/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// corsMiddleware menangani izin akses Cross-Origin Resource Sharing (CORS)
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

func main() {
	// Inisialisasi structured logger (slog) sebagai logger default
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Memuat environment variables dari file .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("Peringatan: Gagal memuat file .env")
	}

	// Memuat konfigurasi pengguna dari config.yaml
	if err := auth.LoadConfig("config.yaml"); err != nil {
		slog.Error("Fatal: Gagal memuat config.yaml", "error", err)
		os.Exit(1)
	}

	// Memuat JWT secret key dari environment
	if err := auth.LoadJWTSecret(); err != nil {
		slog.Error("Fatal: Gagal memuat JWT secret", "error", err)
		os.Exit(1)
	}

	// Membuat router baru menggunakan gorilla/mux
	router := mux.NewRouter()
	router.Use(corsMiddleware) // Terapkan CORS ke semua rute

	// Mendaftarkan rute API
	api := router.PathPrefix("/api").Subrouter()

	// Rute publik untuk otentikasi
	api.HandleFunc("/login", handlers.LoginHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST", "OPTIONS")

	// Subrouter untuk rute film yang diproteksi
	moviesRouter := api.PathPrefix("/movies").Subrouter()
	moviesRouter.Use(middleware.AuthMiddleware) // Terapkan middleware auth hanya ke rute film
	moviesRouter.HandleFunc("", handlers.GetMoviesHandler).Methods("GET")
	moviesRouter.HandleFunc("", handlers.CreateMovieHandler).Methods("POST", "OPTIONS")
	moviesRouter.HandleFunc("/{id}", handlers.GetMovieByIDHandler).Methods("GET")
	moviesRouter.HandleFunc("/{id}", handlers.UpdateMovieHandler).Methods("PUT", "OPTIONS")
	moviesRouter.HandleFunc("/{id}", handlers.DeleteMovieHandler).Methods("DELETE", "OPTIONS")

	// Rute untuk health check
	router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Membaca port dari .env, dengan nilai default "8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := ":" + port
	fmt.Printf("🚀 Server is running on port %s\n", port)
	fmt.Printf("📡 API endpoints available at:\n")
	fmt.Printf("   POST   /api/movies      - Create a new movie\n")
	fmt.Printf("   GET    /api/movies      - Get all movies\n")
	fmt.Printf("   GET    /api/movies/{id} - Get movie by ID\n")
	fmt.Printf("   PUT    /api/movies/{id} - Update movie by ID\n")
	fmt.Printf("   DELETE /api/movies/{id} - Delete movie by ID\n")
	fmt.Printf("   GET    /health          - Health check\n")
	fmt.Printf("🌐 Server URL: http://localhost:%s\n", port)

	// Menjalankan server
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("Gagal menjalankan server", "error", err)
		log.Fatal(err) // log.Fatal akan menghentikan program
	}

}