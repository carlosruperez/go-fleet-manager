package usecases

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/common"
	"github.com/go-fleet-manager/internal/repository"
)

type requestBody struct {
	FullName string `json:"full_name"`
	Url      string `json:"html_url"`
}

type branch struct {
	Name string `json:"name"`
}

type tag struct {
	Name string `json:"name"`
}

func GetRepositories() []repository.Repository {
	outputBytes, _, err := gh.Exec("repo", "list", "carlosruperez", "-L", "100")
	if err != nil {
		panic(err)
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(strings.NewReader(outputBytes.String()))

	// Slice to hold repositories
	var repos []requestBody

	// Iterate over each line of output
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line by whitespace
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			repo := requestBody{
				FullName: parts[0],
				Url:      parts[1],
			}
			repos = append(repos, repo)
		}
	}

	repositories := []repository.Repository{}
	repoType := common.Others
	for _, repo := range repos {

		if strings.Contains(repo.FullName, "-ms-") && strings.Contains(repo.FullName, "-api-server") {
			repoType = common.Microservice
		} else if strings.Contains(repo.FullName, "-sdk") {
			repoType = common.SDK
		} else {
			repoType = common.Others
		}
		repoName := strings.Replace(repo.FullName, "carlosruperez/", "", 1)
		repositories = append(repositories, repository.Repository{
			FullName: repo.FullName,
			Name:     repoName,
			Url:      repo.Url,
			Type:     repoType,
		})
	}
	sort.Slice(repositories, func(i, j int) bool {
		return repositories[i].Name < repositories[j].Name
	})
	return repositories
}

func UpdateMinor(repository repository.Repository, ctx context.Context) (stdout, stderr bytes.Buffer, err error) {

	devMinorConfig := config.GetDevMinorConfig()
	workflowId := devMinorConfig.WorkflowId
	workflowRepository := devMinorConfig.WorkflowRepository

	stdout, stderr, err = gh.ExecContext(ctx, "workflow", "run", workflowId, "-R", workflowRepository, "-F", "version=minor", "-F", "appGithubRepository="+repository.FullName)
	return stdout, stderr, err
}

func GetRepositoryBranches(repository repository.Repository, ctx context.Context) ([]string, error) {
	var branches []branch
	endpoint := "repos/" + repository.FullName + "/branches"
	outputBytes, _, err := gh.ExecContext(ctx, "api", endpoint)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outputBytes.Bytes(), &branches)
	if err != nil {
		return nil, err
	}

	var branchStrs []string

	for _, branch := range branches {
		branchStrs = append(branchStrs, branch.Name)
	}

	return branchStrs, nil
}

func GetRepositoryTags(repository repository.Repository, ctx context.Context) ([]string, error) {
	var tags []tag
	endpoint := "repos/" + repository.FullName + "/tags"
	outputBytes, _, err := gh.ExecContext(ctx, "api", endpoint)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outputBytes.Bytes(), &tags)
	if err != nil {
		return nil, err
	}

	var tagStrs []string

	for _, tag := range tags {
		tagStrs = append(tagStrs, tag.Name)
	}

	return tagStrs, nil
}
