package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/goforbroke1006/localhost-compose/domain"
)

var (
	configFilename    = "localhost-compose.yml"
	schema            domain.ComposeSchema
	composeWorkingDir string
)

var rootCmd = &cobra.Command{
	Use: "localhost-compose",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFilename, "file", "f", configFilename,
		fmt.Sprintf("Specify an alternate compose file\n                              (default: %s)", configFilename))
}

func Execute() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	composeWorkingDir = path

	composeBytes, err := os.ReadFile(configFilename)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(composeBytes, &schema); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
