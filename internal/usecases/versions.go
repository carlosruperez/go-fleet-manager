package usecases

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-fleet-manager/internal/common"
	"github.com/go-fleet-manager/internal/repository"
)

type reqContent struct {
	Version string `json:"version"`
}

// I need a mapping of the environment to the url
var url = "http://carlosruperez.internal-api.dev.internal/version"

var environmentUrls = map[common.Environment]string{
	common.Development: "http://carlosruperez.internal-api.dev.internal/",
	common.Staging:     "http://carlosruperez.internal-api.pre.internal/",
	common.Production:  "http://carlosruperez.internal-api.internal/",
}

func GetVersion(repo repository.Repository, environment common.Environment) (string, error) {

	msPath, err := repo.GetMSPath()
	if err != nil {
		return "", errors.New("Repository is not a microservice")
	}

	fullUrl := environmentUrls[environment] + msPath + "/version"
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
