package usecases

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/common"
	"github.com/go-fleet-manager/internal/repository"
)

type reqContent struct {
	Version string `json:"version"`
}

func GetVersion(repo repository.Repository, environment common.Environment) (string, error) {
	environments := config.GetEnvironments()
	environmentUrl, err := environments.GetEnvironmentUrl(environment)
	if err != nil {
		return "", err
	}

	msPath, err := repo.GetMSPath()
	if err != nil {
		return "", errors.New("Repository is not a microservice")
	}

	fullUrl := environmentUrl + msPath + "/version"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(fullUrl)
	if err != nil {
		return "", errors.New("Error getting version " + fullUrl)
	}
	defer resp.Body.Close()

	content := reqContent{}
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return "", errors.New("Error decoding version" + fullUrl)
	}

	return content.Version, nil
}
