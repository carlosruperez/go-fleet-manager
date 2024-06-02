package usecases

import (
	"bytes"
	"context"

	"github.com/cli/go-gh/v2"
	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/repository"
)

func Deploy(repo repository.Repository, version string, ctx context.Context) (stdout, stderr bytes.Buffer, err error) {
	appName := repo.GetAppName()

	deployConfig := config.GetDeploymentConfig()
	workflowId := deployConfig.WorkflowId
	workflowRepo := deployConfig.WorkflowRepository

	stdout, stderr, err = gh.ExecContext(ctx, "workflow", "run", workflowId, "-R", workflowRepo, "-F", "appName="+appName, "-F", "dockerImageVersion="+version, "-F", "env=prod")
	return stdout, stderr, err
}
