package commands

import (
	"fmt"

	"github.com/abassian/huron/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of huron being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
