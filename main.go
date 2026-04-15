package main

import (
	"database/sql"
	"gohttp/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	const (
		fileRootPath = "."
		port         = "8080"
	)
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()

	// Handler
	fsHandler := cfg.middlewareMetricsInc(http.FileServer(http.Dir(fileRootPath)))
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpsDelete)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/login", cfg.handlerUsersLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeToken)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaUpgrade)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: localhost:%s\n", fileRootPath, port)
	log.Fatal(httpServer.ListenAndServe())
}
