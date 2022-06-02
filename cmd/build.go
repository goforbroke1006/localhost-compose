package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/goforbroke1006/localhost-compose/domain"
	"github.com/goforbroke1006/localhost-compose/internal"
	"github.com/goforbroke1006/localhost-compose/pkg"
)

func NewBuildCmd() *cobra.Command {
	var (
		composeSchema     domain.ComposeSchema
		composeWorkingDir string
	)

	return &cobra.Command{
		Use: "build",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			composeWorkingDir = path

			composeBytes, err := os.ReadFile(composeFilename)
			if err != nil {
				return err
			}
			if err = yaml.Unmarshal(composeBytes, &composeSchema); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := pkg.NewLogger()

			ctx := context.Background()

			for svcName, svcSpec := range composeSchema.Services {

				if len(svcSpec.Build.Shell) == 0 {
					logger.Infof(svcName, "skipped")
					continue
				}

				var currentWorkDir string
				if filepath.IsAbs(svcSpec.WorkingDir) {
					currentWorkDir = svcSpec.WorkingDir
				} else {
					currentWorkDir = filepath.Join(composeWorkingDir, svcSpec.WorkingDir)
				}

				executor := internal.NewBashRunner()
				stdOut := make(chan string)
				stdErr := make(chan string)

				go func() {
				Loop:
					for {
						select {
						case msg, opened := <-stdOut:
							if !opened {
								break Loop
							}
							logger.Info(svcName, msg)
						case msg, opened := <-stdErr:
							if !opened {
								break Loop
							}
							logger.Errorf(svcName, msg)
						}
					}
				}()

				code, reason, err := executor.ExecuteWithContext(ctx, svcSpec.Build.Shell, currentWorkDir, domain.RunnerModeOneShell, stdOut, stdErr)
				if err != nil {
					panic(err)
				}
				time.Sleep(time.Second)
				logger.Infof(svcName, "exit code: %d, reason: %s", code, reason)

			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewBuildCmd())
}
