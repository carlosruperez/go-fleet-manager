//go:build main

package config

import (
	_ "embed"
)

//go:embed config-main.yaml
var configFile string
