package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const (
		fileRootPath = "."
		port         = "8080"
	)

	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	// Handler
	fsHandler := cfg.middlewareMetricsInc(http.FileServer(http.Dir(fileRootPath)))
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))
	
	mux.HandleFunc("GET /api/healthz/", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: localhost:%s\n", fileRootPath, port)
	log.Fatal(httpServer.ListenAndServe())
}

type apiConfig struct {
	fileServerHits atomic.Int32
}
