package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	"gopkg.in/yaml.v3"
)

func main() {
	composeFilename := "localhost-compose.yml"

	if len(os.Args) > 1 {
		composeFilename = os.Args[1]
	}

	composeBytes, err := os.ReadFile(composeFilename)
	if err != nil {
		panic(err)
	}
	var schema ComposeSchema
	if err = yaml.Unmarshal(composeBytes, &schema); err != nil {
		panic(err)
	}

	//fmt.Println(schema)

	ctx, cancel := context.WithCancel(context.Background())
	servicesWg := sync.WaitGroup{}

	for svcName, svcSpec := range schema.Services {
		servicesWg.Add(1)
		go func(svcName string, svcSpec ServiceSpec) {
			{
				buildCmd := exec.Command("/bin/sh", "-c", svcSpec.Build.Shell)
				buildCmd.Stdout = os.Stdout
				if err := buildCmd.Run(); err != nil {
					fmt.Println("ERROR:", err.Error(), ":", "build", svcName)
					return
				}
				//output, err := buildCmd.Output()
				//if err != nil {
				//	fmt.Println("ERROR:", err.Error(), ":", "build", svcName)
				//	return
				//}
				//fmt.Println("INFO:", svcName, "build output", string(output))
				fmt.Println("INFO:", svcName, "build exit", buildCmd.ProcessState.ExitCode())
			}

			{
				commandCmd := exec.CommandContext(ctx, "/bin/sh", "-c", svcSpec.Command)
				//commandCmd := exec.Command("/bin/sh", "-c", svcSpec.Command)
				commandCmd.Stdout = os.Stdout
				go func() {
					if err := commandCmd.Run(); err != nil {
						fmt.Println("ERROR:", err.Error(), ":", svcName, svcSpec.Command)
					}
					fmt.Println("INFO:", svcName, "command exit", commandCmd.ProcessState.ExitCode())
				}()
				fmt.Println("INFO:", svcName, "running")
			}

			<-ctx.Done()
			fmt.Println("INFO:", svcName, "stopping")
			servicesWg.Done()

		}(svcName, svcSpec)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Stopping all")
	cancel()
	servicesWg.Wait()
}

type ComposeSchema struct {
	Services map[string]ServiceSpec `yaml:"services"`
}

type ServiceSpec struct {
	WorkingDir string    `yaml:"working_dir"`
	Build      BuildSpec `yaml:"build"`
	Command    string    `yaml:"command"`
}

type BuildSpec struct {
	Shell string `yaml:"shell"`
}
