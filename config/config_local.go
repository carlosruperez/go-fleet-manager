//go:build local

package config

import (
	_ "embed"
)

//go:embed config-local.yaml
var configFile string
