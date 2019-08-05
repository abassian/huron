package commands

import (
	"github.com/spf13/cobra"
)

var (
	config = NewDefaultCLIConfig()
)

//RootCmd is the root command for Huron
var RootCmd = &cobra.Command{
	Use:              "huron",
	Short:            "huron consensus",
	TraverseChildren: true,
}
