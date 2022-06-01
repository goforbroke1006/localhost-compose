package domain

import "context"

type RunnerMode string

const (
	RunnerModeSplit    = "split"
	RunnerModeOneShell = "one-shell"
)

type RunnerExitReason string

const (
	RunnerExitReasonUnknown = "unknown"
	RunnerExitReasonDone    = "done"
	RunnerExitReasonKilled  = "killed"
)

type CommandRunner interface {
	ExecuteWithContext(
		ctx context.Context,
		commands string,
		workingDir string,
		mode RunnerMode,
		stdOut, stdErr chan string,
	) (code int, reason RunnerExitReason, err error)
}
