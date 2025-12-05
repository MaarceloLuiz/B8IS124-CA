package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MaarceloLuiz/worldle-replica/pkg/game"
	"github.com/MaarceloLuiz/worldle-replica/pkg/geography/geocalc"
	terr "github.com/MaarceloLuiz/worldle-replica/pkg/geography/territories"
)

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		sessionID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	err := game.StartNewGame(sessionID)
	if err != nil {
		http.Error(w, "Failed to start new game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"sessionId": sessionID,
		"message":   "New game started",
	})
}

func SilhouetteHandler(w http.ResponseWriter, r *http.Request) {
	silhouette, err := game.GetCurrentSilhouette()
	if err != nil {
		http.Error(w, "Failed to fetch silhouette", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(silhouette)
}

func AllTerritoriesHandler(w http.ResponseWriter, r *http.Request) {
	territories, err := terr.GetFormattedTerritoryNames()
	if err != nil {
		http.Error(w, "Failed to fetch territories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(territories)
	if err != nil {
		http.Error(w, "Failed to encode territories", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func AnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Use GET to retrieve the answer."}`))
		return
	}

	answerCountry := game.State.Country
	if answerCountry == "" {
		http.Error(w, "Game not initialized", http.StatusInternalServerError)
		return
	}

	mapsURL, err := geocalc.GetMapsURLAnswer(answerCountry)
	if err != nil {
		http.Error(w, "Failed to generate maps URL for the answer", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"answer": answerCountry,
		"url":    mapsURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GuessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Use POST with a JSON payload to make a guess."}`))
		return
	}

	var guessCountry struct {
		Guess string `json:"guess"`
	}
	if err := json.NewDecoder(r.Body).Decode(&guessCountry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	guessCountry.Guess = strings.ReplaceAll(guessCountry.Guess, " ", "_")

	answerCountry := game.State.Country
	if answerCountry == "" {
		http.Error(w, "Game not initialized", http.StatusInternalServerError)
		return
	}

	// Case-insensitive comparison of guess and current country.
	isCorrect := strings.EqualFold(guessCountry.Guess, answerCountry)

	distance, direction, err := geocalc.GetDistanceAndDirection(guessCountry.Guess, answerCountry)
	if err != nil {
		http.Error(w, "Failed to calculate distance", http.StatusInternalServerError)
		return
	}

	var mapsURL string
	if isCorrect {
		mapsURL, err = geocalc.GetMapsURL(answerCountry, answerCountry)
		direction = ""
	} else {
		mapsURL, err = geocalc.GetMapsURL(guessCountry.Guess, answerCountry)
	}

	if err != nil {
		http.Error(w, "Failed to generate maps URL", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"isCorrect": isCorrect,
		"distance":  distance,
		"direction": direction,
		"url":       mapsURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
