package repository

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/config"
)

type commandRepository struct {
	config *config.Command
}

func NewCommandRepository(conf *config.Command) api.Repository {
	return &commandRepository{config: conf}
}

func (r *commandRepository) Exec(command string) (output string, exitCode int, duration time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(r.config.Timeout))
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err = cmd.Run()
	duration = time.Since(start)
	output = buf.String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			err = nil
		} else if ctx.Err() == context.DeadlineExceeded {
			exitCode = -1
			err = ctx.Err()
		}
	} else if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return
}
