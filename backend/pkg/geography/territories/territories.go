package geography

import (
	"strings"
	"sync"

	firebase "github.com/MaarceloLuiz/worldle-replica/pkg/firebase"
	"github.com/sirupsen/logrus"
)

var (
	territoriesCache []string
	cacheMutex       sync.RWMutex
	cacheInitialized bool
)

func GetFormattedTerritoryNames() ([]string, error) {
	territories, err := GetAllTerritories()
	if err != nil {
		logrus.Error("Failed to get all territories")
		return nil, err
	}

	formatted := make([]string, len(territories))
	for i, territory := range territories {
		formatted[i] = strings.ReplaceAll(strings.ToUpper(territory), "_", " ")
	}

	return formatted, nil
}

func GetAllTerritories() ([]string, error) {
	// cache the result to avoid multiple calls to Firebase
	cacheMutex.RLock()
	if cacheInitialized {
		logrus.Info("Returning cached territories")
		defer cacheMutex.RUnlock()
		return territoriesCache, nil
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cacheInitialized {
		logrus.Info("Returning cached territories")
		return territoriesCache, nil
	}

	logrus.Info("Loading territories from Firebase")
	territories, err := firebase.GetAllTerritoriesFromStorage()
	if err != nil {
		logrus.Errorf("Failed to load territories from Firebase: %v", err)
		return nil, err
	}

	territoriesCache = territories
	cacheInitialized = true
	logrus.Infof("Cached %d territories", len(territories))

	return territories, nil
}
