package silhouettes

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	geography "github.com/MaarceloLuiz/worldle-replica/pkg/geography/territories"
	"github.com/MaarceloLuiz/worldle-replica/pkg/github"
	"github.com/sirupsen/logrus"
)

func FetchSilhouette(country string) ([]byte, error) {
	endpoint := fmt.Sprintf("contents/silhouettes/%s.png", country)
	headers := map[string]string{
		"Accept": "application/vnd.github.raw",
	}

	request, err := github.CreateGitHubRequest(endpoint, headers)
	if err != nil {
		logrus.Error("Error creating new GitHub request")
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Errorf("Failed to make request to GitHub API - Authorization issue or network error: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch silhouette from GitHub API")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body from GitHub API: %v", err)
		return nil, err
	}

	return data, nil
}

func GetRandomCountry() (string, error) {
	countries, err := geography.GetAllTerritories()
	if err != nil {
		logrus.Error("Failed to get all territories")
		return "", err
	}

	seed := time.Now().UnixNano() // to avoid seeding the same number every time
	random := rand.New(rand.NewSource(seed))

	randomIndex := random.Intn(len(countries))
	randomCountry := countries[randomIndex]

	return randomCountry, nil
}
