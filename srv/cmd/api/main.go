package main

import (
	"context"
	"log"
	"net/http"

	"regres/srv/internal/config"
	"regres/srv/internal/database"
	"regres/srv/internal/database/queries"
	"regres/srv/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.Load()

	ctx := context.Background()

	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := queries.New(pool)

	mux := http.NewServeMux()
	routes.Register(mux, q)

	handler := withCORS(mux)

	port := config.Port()

	log.Println(config.AppName() + " API running on port " + port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == config.CorsAllowedOrigin() {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
