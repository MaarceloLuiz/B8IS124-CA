package geocalc

import (
	"fmt"
	"math"
	"sync"
)

var (
	coordinatesCache = make(map[string][2]float64)
	cacheMutex       sync.RWMutex
)

func GetDistanceAndDirection(guess string, answer string) (float64, string, error) {
	guessLat, guessLng, err := getCachedCoordinates(guess)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get coordinates for guess: %w", err)
	}

	answerLat, answerLng, err := getCachedCoordinates(answer)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get coordinates for answer: %w", err)
	}

	distance, err := haversine(guessLat, guessLng, answerLat, answerLng)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate haversine distance: %w", err)
	}

	distanceKM := distance / 1000
	distanceKM = math.Round(distanceKM)
	direction := getCompass(guessLat, guessLng, answerLat, answerLng)

	return distanceKM, direction, nil
}

func GetMapsURL(guess string, answer string) (string, error) {
	guessLat, guessLng, err := getCachedCoordinates(guess)
	if err != nil {
		return "", fmt.Errorf("failed to get coordinates for guess: %w", err)
	}

	answerLat, answerLng, err := getCachedCoordinates(answer)
	if err != nil {
		return "", fmt.Errorf("failed to get coordinates for answer: %w", err)
	}

	url := fmt.Sprintf("https://www.google.com/maps/dir/%f,%f/%f,%f",
		guessLat, guessLng, answerLat, answerLng)

	return url, nil
}

func GetMapsURLAnswer(answer string) (string, error) {
	url := fmt.Sprintf("https://www.google.com/maps/place/%s", answer)

	return url, nil
}

func getCachedCoordinates(territory string) (float64, float64, error) {
	cacheMutex.RLock()
	coords, found := coordinatesCache[territory]
	cacheMutex.RUnlock()

	if found {
		return coords[0], coords[1], nil
	}

	lat, lng, err := geocode(territory)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to geocode territory: %s: %w", territory, err)
	}

	cacheMutex.Lock()
	coordinatesCache[territory] = [2]float64{lat, lng}
	cacheMutex.Unlock()

	return lat, lng, nil
}
