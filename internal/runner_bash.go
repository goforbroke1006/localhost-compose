package internal

import (
	"bytes"
	"context"
	"errors"
	"os/exec"

	"github.com/goforbroke1006/localhost-compose/domain"
	"github.com/goforbroke1006/localhost-compose/pkg"
)

func NewBashRunner() *bashRunner {
	return &bashRunner{}
}

var _ domain.CommandRunner = &bashRunner{}

type bashRunner struct{}

func (r bashRunner) ExecuteWithContext(
	ctx context.Context,
	commands string,
	workingDir string,
	mode domain.RunnerMode,
	stdOut, stdErr chan string,
) (code int, reason domain.RunnerExitReason, err error) {
	if mode == domain.RunnerModeOneShell {
		return r.exec(ctx, commands, workingDir, stdOut, stdErr)
	}

	if mode == domain.RunnerModeSplit {
		list := r.splitCommands(commands)
		for _, command := range list {
			code, reason, err = r.exec(ctx, command, workingDir, stdOut, stdErr)
			if err != nil {
				return code, reason, err
			}
			if reason != domain.RunnerExitReasonDone {
				return code, reason, err
			}
		}
	}

	return 0, "", errors.New("wrong runner mode: " + string(mode))
}

func (r bashRunner) exec(ctx context.Context,
	commands string,
	workingDir string,
	stdOut, stdErr chan string,
) (code int, reason domain.RunnerExitReason, err error) {
	commandCmd := exec.CommandContext(ctx, "/bin/bash", "-c", commands)
	commandCmd.Dir = workingDir

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	commandCmd.Stdout = stdout // standard output
	commandCmd.Stderr = stderr // standard error

	readOut := pkg.NewBashOutputReader(stdout)
	readErr := pkg.NewBashOutputReader(stderr)

	if err := commandCmd.Start(); err != nil {
		return 0, "", err
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				length, text, _ := readOut.ReadString()
				if length == 0 {
					continue
				}
				stdOut <- text
			}
		}
	}(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				length, text, _ := readErr.ReadString()
				if length == 0 {
					continue
				}
				stdErr <- text
			}
		}
	}()

	if err := commandCmd.Wait(); err != nil {
		if err.Error() == "signal: interrupt" || err.Error() == "signal: killed" {
			reason = domain.RunnerExitReasonKilled
		} else {
			reason = domain.RunnerExitReasonUnknown
		}
		return code, reason, err
	}

	code = commandCmd.ProcessState.ExitCode()
	reason = domain.RunnerExitReasonDone

	return code, reason, nil
}

func (r bashRunner) splitCommands(commands string) []string {
	panic("implement me")
}
