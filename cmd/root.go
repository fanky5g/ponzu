package cmd

import (
	"fmt"
	"github.com/fanky5g/ponzu/cmd/generate"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:  "ponzu",
	Long: `Ponzu is an open-source HTTP server framework and CMS`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(RootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}
}

func init() {
	generate.RegisterCommandRecursive(RootCmd)
}
