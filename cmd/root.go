package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	composeFilename = "localhost-compose.yml"
)

var rootCmd = &cobra.Command{
	Use: "localhost-compose",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&composeFilename, "file", "f", composeFilename,
		fmt.Sprintf("Specify an alternate compose file\n                              (default: %s)", composeFilename))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
