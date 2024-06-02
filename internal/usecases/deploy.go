package usecases

import (
	"bytes"
	"context"
	"github.com/go-fleet-manager/internal/common"

	"github.com/cli/go-gh/v2"
	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/repository"
)

func Deploy(repo repository.Repository, version string, environment common.Environment, ctx context.Context) (stdout, stderr bytes.Buffer, err error) {
	deployConfig := config.GetDeploymentConfig()
	workflowId := deployConfig.WorkflowId
	workflowRepo := deployConfig.WorkflowRepository

	var env string
	switch environment {
	case common.Development:
		env = "dev"
	case common.Staging:
		env = "pre"
	case common.Production:
		env = "prod"
	}

	stdout, stderr, err = gh.ExecContext(ctx, "workflow", "run", workflowId, "-R", workflowRepo, "-F", "appName="+repo.Name, "-F", "dockerImageVersion="+version, "-F", "env="+env)
	return stdout, stderr, err
}
