package github

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type GitConfig struct {
	Token  string `env:"GITHUB_TOKEN"`
	Repo   string `env:"GITHUB_REPO"`
	Owner  string `env:"GITHUB_OWNER"`
	Branch string `env:"GITHUB_BRANCH"`
}

func CreateGitHubRequest(endpoint string, headers map[string]string) (*http.Request, error) {
	gitConfig, err := loadGitConfig()
	if err != nil {
		logrus.Error("Error loading GitHub configuration")
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/%s?ref=%s",
		gitConfig.Owner, gitConfig.Repo, endpoint, gitConfig.Branch)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request to GitHub API")
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+gitConfig.Token)
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	return request, nil
}

func loadGitConfig() (*GitConfig, error) {
	gitConfig := &GitConfig{
		Token:  os.Getenv("GITHUB_TOKEN"),
		Repo:   os.Getenv("GITHUB_REPO"),
		Owner:  os.Getenv("GITHUB_OWNER"),
		Branch: os.Getenv("GITHUB_BRANCH"),
	}

	if gitConfig.Token == "" || gitConfig.Repo == "" || gitConfig.Owner == "" || gitConfig.Branch == "" {
		return nil, fmt.Errorf("missing required GitHub configuration")
	}

	return gitConfig, nil
}
