package config

import (
	"errors"
	"log"

	"github.com/go-fleet-manager/internal/common"
	"gopkg.in/yaml.v2"
)

type CacheConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type WorkflowConfig struct {
	WorkflowId         string `yaml:"workflowId"`
	WorkflowRepository string `yaml:"workflowRepository"`
}

type EnvironmentConfig struct {
	Development string `yaml:"development"`
	Staging     string `yaml:"staging"`
	Production  string `yaml:"production"`
}

func (e EnvironmentConfig) GetEnvironmentUrl(environment common.Environment) (string, error) {
	switch environment {
	case common.Development:
		return e.Development, nil
	case common.Staging:
		return e.Staging, nil
	case common.Production:
		return e.Production, nil
	}
	return "", errors.New("Environment not found")
}

type Config struct {
	Cache              []CacheConfig     `yaml:"cache"`
	DevMinor           WorkflowConfig    `yaml:"devMinor"`
	MaintenanceMode    WorkflowConfig    `yaml:"maintenanceMode"`
	Deployment         WorkflowConfig    `yaml:"deployment"`
	CreateMainPRConfig WorkflowConfig    `yaml:"createMainPR"`
	Environments       EnvironmentConfig `yaml:"environments"`
}

func getConfig() Config {
	var config Config
	err := yaml.Unmarshal([]byte(configFile), &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}

func GetCaches() []CacheConfig {
	config := getConfig()

	cacheConfigs := []CacheConfig{}
	cacheConfigs = append(cacheConfigs, config.Cache...)

	return cacheConfigs
}

func GetDevMinorConfig() WorkflowConfig {
	config := getConfig()

	return config.DevMinor
}

func GetMaintenanceModeConfig() WorkflowConfig {
	config := getConfig()

	return config.MaintenanceMode
}

func GetDeploymentConfig() WorkflowConfig {
	config := getConfig()

	return config.Deployment
}

func GetCreateMainPRConfig() WorkflowConfig {
	config := getConfig()

	return config.CreateMainPRConfig
}

func GetEnvironments() EnvironmentConfig {
	config := getConfig()

	return config.Environments
}
