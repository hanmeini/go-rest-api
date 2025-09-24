package main

import (
	"fmt"
	"log"
	"net/http"

	"go-flix-api/handlers"

	"github.com/gorilla/mux"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Mengizinkan semua origin
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
	// Create new router using gorilla/mux
	router := mux.NewRouter()

	// Register API routes
	api := router.PathPrefix("/api").Subrouter()

	// Movies routes
	api.HandleFunc("/movies", handlers.CreateMovieHandler).Methods("POST", "OPTIONS") 
	api.HandleFunc("/movies", handlers.GetMoviesHandler).Methods("GET")
	api.HandleFunc("/movies/{id}", handlers.GetMovieByIDHandler).Methods("GET")
	api.HandleFunc("/movies/{id}", handlers.UpdateMovieHandler).Methods("PUT", "OPTIONS") 
	api.HandleFunc("/movies/{id}", handlers.DeleteMovieHandler).Methods("DELETE", "OPTIONS")

	// Health check route
	router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Apply CORS middleware to router
	// Cara Anda menggunakan .Use() sudah benar dan ini adalah cara yang lebih modern!
	router.Use(corsMiddleware)

	// Start server on port 8080
	port := ":8080"
	fmt.Printf("🚀 Server is running on port %s\n", port)
	fmt.Printf("📡 API endpoints available at:\n")
	fmt.Printf("   POST   /api/movies      - Create a new movie\n")
	fmt.Printf("   GET    /api/movies      - Get all movies\n")
	fmt.Printf("   GET    /api/movies/{id} - Get movie by ID\n")
	fmt.Printf("   PUT    /api/movies/{id} - Update movie by ID\n")
	fmt.Printf("   DELETE /api/movies/{id} - Delete movie by ID\n")
	fmt.Printf("   GET    /health          - Health check\n")
	fmt.Printf("🌐 Server URL: http://localhost%s\n", port)

	// Start the server
	log.Fatal(http.ListenAndServe(port, router))
}
