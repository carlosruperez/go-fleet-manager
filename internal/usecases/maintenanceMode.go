package usecases

import (
	"bytes"
	"context"

	"github.com/cli/go-gh/v2"
	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/common"
)

type MaintenanceAction string

const (
	EnableMaintenance  MaintenanceAction = "enable"
	DisableMaintenance MaintenanceAction = "disable"
)

func MaintenanceMode(action MaintenanceAction, environment common.Environment, ctx context.Context) (stdout, stderr bytes.Buffer, err error) {
	maintenanceModeConfig := config.GetMaintenanceModeConfig()
	workflowId := maintenanceModeConfig.WorkflowId
	workflowRepo := maintenanceModeConfig.WorkflowRepository

	mapEnvironments := map[common.Environment]string{
		common.Development: "dev",
		common.Staging:     "preprod",
		common.Production:  "prod",
	}

	environmentName := mapEnvironments[environment]

	stdout, stderr, err = gh.ExecContext(ctx, "workflow", "run", workflowId, "-R", workflowRepo, "-F", "action="+string(action), "-F", "environment="+environmentName)
	return stdout, stderr, err
}
