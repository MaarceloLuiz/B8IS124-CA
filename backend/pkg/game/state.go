package game

import (
	"errors"
	"strings"
	"sync"

	"github.com/MaarceloLuiz/worldle-replica/pkg/geography/silhouettes"
	"github.com/sirupsen/logrus"
)

var ErrGameNotInitialized = errors.New("game not initialized")

type GameState struct {
	Country   string
	Image     []byte
	SessionID string
	Mutex     sync.RWMutex
}

var State = GameState{}

func StartNewGame(sessionID string) error {
	State.Mutex.Lock()
	defer State.Mutex.Unlock()

	if State.SessionID == sessionID && State.Country != "" {
		// logrus.Infof("Reusing existing game for session %s: %s", sessionID, State.Country) // DEBUG LOG
		logrus.Info("Reusing existing game session")
		return nil
	}

	country, err := silhouettes.GetRandomCountry()
	if err != nil {
		return err
	}

	img, err := silhouettes.FetchSilhouette(country)
	if err != nil {
		return err
	}

	State.Country = strings.ReplaceAll(strings.ToUpper(country), "_", " ")
	State.Image = img
	State.SessionID = sessionID
	return nil
}

func GetCurrentSilhouette() ([]byte, error) {
	State.Mutex.RLock()
	defer State.Mutex.RUnlock()

	if State.Image == nil {
		return nil, ErrGameNotInitialized
	}

	return State.Image, nil
}
