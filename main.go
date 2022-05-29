package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"localhost-compose/domain"
	"localhost-compose/pkg"
)

func main() {
	logger := pkg.NewLogger()

	composeFilename := "localhost-compose.yml"

	if len(os.Args) > 1 {
		composeFilename = os.Args[1]
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	composeBytes, err := os.ReadFile(composeFilename)
	if err != nil {
		panic(err)
	}
	var schema domain.ComposeSchema
	if err = yaml.Unmarshal(composeBytes, &schema); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	servicesWg := sync.WaitGroup{}

	for svcName, svcSpec := range schema.Services {
		servicesWg.Add(1)
		go func(svcName string, svcSpec domain.ServiceSpec) {
			//{
			//	buildCmd := exec.Command("/bin/sh", "-c", svcSpec.Build.Shell)
			//	buildCmd.Stdout = os.Stdout
			//	if err := buildCmd.Run(); err != nil {
			//		fmt.Println("ERROR:", err.Error(), ":", "build", svcName)
			//		return
			//	}
			//
			//	fmt.Println("INFO:", svcName, "build exit", buildCmd.ProcessState.ExitCode())
			//}

			var currentWorkDir string
			if filepath.IsAbs(svcSpec.WorkingDir) {
				currentWorkDir = svcSpec.WorkingDir
			} else {
				currentWorkDir = filepath.Join(path, svcSpec.WorkingDir)
			}

			{
				commandCmd := exec.CommandContext(ctx, "/bin/bash", "-c", svcSpec.Command)
				commandCmd.Dir = currentWorkDir

				stdout := new(bytes.Buffer)
				stderr := new(bytes.Buffer)

				commandCmd.Stdout = stdout // standard output
				commandCmd.Stderr = stderr // standard error

				readOut := pkg.NewBashOutputReader(stdout)
				readErr := pkg.NewBashOutputReader(stderr)

				if err = commandCmd.Start(); err != nil {
					panic(err)
				}

				go func() {
					for {
						length, text, _ := readOut.ReadString()
						if length == 0 {
							continue
						}

						logger.Info(svcName, text)

					}
				}()
				go func() {
					for {
						length, text, _ := readErr.ReadString()
						if length == 0 {
							continue
						}

						logger.Info(svcName, text)
					}
				}()

				if err := commandCmd.Wait(); err != nil {
					panic(err)
				}

				if commandCmd.ProcessState.ExitCode() == 0 {
					logger.Infof(svcName, "exit code %d", commandCmd.ProcessState.ExitCode())
				} else {
					logger.Errorf(svcName, "exit code %d", commandCmd.ProcessState.ExitCode())
				}
			}

			logger.Infof(svcName, "stopping")
			servicesWg.Done()

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
}
