package geography

import (
	"strings"

	firebase "github.com/MaarceloLuiz/worldle-replica/pkg/firebase"
	"github.com/sirupsen/logrus"
)

func GetFormattedTerritoryNames() ([]string, error) {
	territories, err := firebase.GetAllTerritories()
	if err != nil {
		logrus.Error("Failed to get all territories")
		return nil, err
	}

	for i, territory := range territories {
		territories[i] = strings.ReplaceAll(strings.ToUpper(territory), "_", " ")
	}

	return territories, nil
}
