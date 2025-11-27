package geography

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MaarceloLuiz/worldle-replica/pkg/github"
	"github.com/sirupsen/logrus"
)

func GetFormattedTerritoryNames() ([]string, error) {
	territories, err := GetAllTerritories()
	if err != nil {
		logrus.Error("Failed to get all territories")
		return nil, err
	}

	for i, territory := range territories {
		territories[i] = strings.ReplaceAll(strings.ToUpper(territory), "_", " ")
	}

	return territories, nil
}

func GetAllTerritories() ([]string, error) {
	request, err := github.CreateGitHubRequest("contents/silhouettes", nil)
	if err != nil {
		logrus.Error("Error creating new GitHub request")
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("Failed to make request to GitHub API - Authorization issue or network error: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch silhouettes from GitHub API: %s", response.Status)

	}

	var files []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(response.Body).Decode(&files); err != nil {
		logrus.Errorf("Failed to decode response body from GitHub API: %v", err)
		return nil, err

	}

	var territories []string
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".png") {
			territories = append(territories, strings.TrimSuffix(file.Name, ".png"))
		}
	}

	if len(territories) == 0 {
		return nil, fmt.Errorf("no territories found in the response")
	}

	return territories, nil
}
