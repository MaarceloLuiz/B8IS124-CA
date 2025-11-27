package main

import (
	"net/http"

	"github.com/MaarceloLuiz/worldle-replica/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/api/newgame", corsMiddleware(api.NewGameHandler))
	http.HandleFunc("/api/silhouette", corsMiddleware(api.SilhouetteHandler))
	http.HandleFunc("/api/territories", corsMiddleware(api.AllTerritoriesHandler))
	http.HandleFunc("/api/answer", corsMiddleware(api.AnswerHandler))
	http.HandleFunc("/api/guess", corsMiddleware(api.GuessHandler))

	logrus.Info("Starting server on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // React dev server
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
