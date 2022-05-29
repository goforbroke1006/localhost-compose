package domain

import (
	"bufio"
	"strings"
)

type ComposeSchema struct {
	Services map[string]ServiceSpec `yaml:"services"`
}

type ServiceSpec struct {
	WorkingDir string    `yaml:"working_dir"`
	Build      BuildSpec `yaml:"build"`
	Command    string    `yaml:"command"`
}

func (ss ServiceSpec) CommandList() []string {
	if len(ss.Command) == 0 {
		return nil
	}

	var result []string

	reader := bufio.NewReader(strings.NewReader(ss.Command))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		strLine := strings.TrimSpace(string(line))

		if len(strLine) == 0 {
			continue
		}

		result = append(result, strLine)
	}

	return result
}

type BuildSpec struct {
	Shell string `yaml:"shell"`
}
