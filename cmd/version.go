package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示当前版本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
		Version: %s
		flag: %v
		`,
			version, flag)
		os.Exit(0)
	},
}
