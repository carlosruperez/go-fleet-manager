package usecases

import (
	"bytes"
	"context"

	"github.com/cli/go-gh/v2"
	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/repository"
)

func CreateMainPR(repo repository.Repository, branch string, ctx context.Context) (stdout, stderr bytes.Buffer, err error) {

	createMainPRConfig := config.GetCreateMainPRConfig()
	workflowId := createMainPRConfig.WorkflowId
	workflowRepo := createMainPRConfig.WorkflowRepository

	stdout, stderr, err = gh.ExecContext(ctx, "workflow", "run", workflowId, "-R", workflowRepo, "-F", "releaseBranch="+branch, "-F", "appGithubRepository="+repo.FullName)
	return stdout, stderr, err
}
