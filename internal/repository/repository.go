package repository

import (
	"errors"
	"strings"

	"github.com/go-fleet-manager/internal/common"
)

type Repository struct {
	FullName string
	Name     string
	Url      string
	Type     common.RepositoryTypes
}

func (r Repository) GetAppName() string {
	// Replace the -ms- and -api-server from r.Name with empty string
	appName := strings.Replace(r.Name, "-ms", "", 1)
	appName = strings.Replace(appName, "-api-server", "", 1)
	return appName
}

func (r Repository) GetMSPath() (string, error) {
	if r.Type != common.Microservice {
		return "", errors.New("Repository is not a microservice")
	}

	msPath := strings.Join(
		[]string{
			strings.Split(r.Name, "-")[0],
			strings.Split(r.Name, "-")[3],
			strings.Split(r.Name, "-")[2],
		},
		"/",
	)

	return msPath, nil
}
