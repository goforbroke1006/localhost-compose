package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/goforbroke1006/localhost-compose/domain"
	"github.com/goforbroke1006/localhost-compose/internal"
	"github.com/goforbroke1006/localhost-compose/pkg"
)

func NewUpCmd() *cobra.Command {
	var (
		composeSchema     domain.ComposeSchema
		composeWorkingDir string
	)

	return &cobra.Command{
		Use: "up",
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

			ctx, cancel := context.WithCancel(context.Background())
			servicesWg := sync.WaitGroup{}

			for svcName, svcSpec := range composeSchema.Services {
				servicesWg.Add(1)
				go func(svcName string, svcSpec domain.ServiceSpec) {
					defer func() {
						servicesWg.Done()
					}()

					var currentWorkDir string
					if filepath.IsAbs(svcSpec.WorkingDir) {
						currentWorkDir = svcSpec.WorkingDir
					} else {
						currentWorkDir = filepath.Join(composeWorkingDir, svcSpec.WorkingDir)
					}

					{
						commandExecutor := internal.NewBashRunner()
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
						code, reason, err := commandExecutor.ExecuteWithContext(ctx, svcSpec.Command, currentWorkDir, domain.RunnerModeOneShell, stdOut, stdErr)
						if err != nil {
							logger.Errorf(svcName, "%v", err)
							return
						}
						time.Sleep(time.Second)
						logger.Infof(svcName, "exit code: %d, reason: %s", code, reason)
					}

				}(svcName, svcSpec)
			}

			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				<-c
				cancel()
				fmt.Println("Stopping all")
			}()

			servicesWg.Wait()
			cancel()
		},
	}
}

func init() {
	rootCmd.AddCommand(NewUpCmd())
}
