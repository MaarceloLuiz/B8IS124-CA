package silhouettes

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	firebase "github.com/MaarceloLuiz/worldle-replica/pkg/firebase"
	"github.com/sirupsen/logrus"
)

func FetchSilhouette(country string) ([]byte, error) {
	imageUrl, err := firebase.GetImageURL(country)
	if err != nil {
		logrus.Errorf("Failed to get image URL from Firebase: %v", err)
		return nil, err
	}

	response, err := http.Get(imageUrl)
	if err != nil {
		logrus.Errorf("Failed to fetch image from Firebase storage: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch silhouette from Firebase storage: status code %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body: %v", err)
		return nil, err
	}

	return data, nil
}

func GetRandomCountry() (string, error) {
	territories, err := firebase.GetAllTerritories()
	if err != nil {
		logrus.Error("Failed to get all territories")
		return "", err
	}

	seed := time.Now().UnixNano() // to avoid seeding the same number every time
	random := rand.New(rand.NewSource(seed))

	randomIndex := random.Intn(len(territories))
	randomTerritory := territories[randomIndex]

	return randomTerritory, nil
}
