package main

import (
	"net/http"
	"os"

	"github.com/MaarceloLuiz/worldle-replica/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/api/newgame", corsMiddleware(api.NewGameHandler))
	http.HandleFunc("/api/silhouette", corsMiddleware(api.SilhouetteHandler))
	http.HandleFunc("/api/territories", corsMiddleware(api.AllTerritoriesHandler))
	http.HandleFunc("/api/answer", corsMiddleware(api.AnswerHandler))
	http.HandleFunc("/api/guess", corsMiddleware(api.GuessHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified (local development)
	}

	logrus.Infof("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000" // Default for local development (react dev server)
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
