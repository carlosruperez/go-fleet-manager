package config

import (
	"log"

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
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Config struct {
	Cache              []CacheConfig       `yaml:"cache"`
	DevMinor           WorkflowConfig      `yaml:"devMinor"`
	MaintenanceMode    WorkflowConfig      `yaml:"maintenanceMode"`
	ProdDeployment     WorkflowConfig      `yaml:"prodDeployment"`
	CreateMainPRConfig WorkflowConfig      `yaml:"createMainPR"`
	Environments       []EnvironmentConfig `yaml:"environments"`
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

func GetProdDeploymentConfig() WorkflowConfig {
	config := getConfig()

	return config.ProdDeployment
}

func GetCreateMainPRConfig() WorkflowConfig {
	config := getConfig()

	return config.CreateMainPRConfig
}

func GetEnvironments() []EnvironmentConfig {
	config := getConfig()

	environmentsConfigs := []EnvironmentConfig{}
	environmentsConfigs = append(environmentsConfigs, config.Environments...)

	return environmentsConfigs
}
