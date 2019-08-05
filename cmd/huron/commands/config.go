package commands

import "github.com/abassian/huron/src/huron"

//CLIConfig contains configuration for the Run command
type CLIConfig struct {
	Huron     huron.HuronConfig `mapstructure:",squash"`
	ProxyAddr  string              `mapstructure:"proxy-listen"`
	ClientAddr string              `mapstructure:"client-connect"`
	Standalone bool                `mapstructure:"standalone"`
}

//NewDefaultCLIConfig creates a CLIConfig with default values
func NewDefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		Huron:     *huron.NewDefaultConfig(),
		ProxyAddr:  "127.0.0.1:1338",
		ClientAddr: "127.0.0.1:1339",
		Standalone: false,
	}
}
