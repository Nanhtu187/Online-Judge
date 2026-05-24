package executor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/docker"
)

type Executor interface {
	Run(ctx context.Context, lang string, binary []byte, input string, memoryLimit int64, timeLimitNano int64) (string, error)
}

type DockerExecutor struct {
	adapter docker.Adapter
}

func NewDockerExecutor(adapter docker.Adapter) *DockerExecutor {
	return &DockerExecutor{adapter: adapter}
}

type config struct {
	image    string
	fileName string
	cmd      []string
}

var langConfigs = map[string]config{
	"go": {
		image:    "alpine:latest",
		fileName: "program",
		cmd:      []string{"./program"},
	},
	"cpp": {
		image:    "alpine:latest",
		fileName: "program",
		cmd:      []string{"./program"},
	},
}

func (e *DockerExecutor) Run(ctx context.Context, lang string, binary []byte, input string, memoryLimit int64, timeLimitNano int64) (string, error) {
	cfg, ok := langConfigs[strings.ToLower(lang)]
	if !ok {
		return "", fmt.Errorf("unsupported language for execution: %s", lang)
	}

	resp, err := e.adapter.RunTask(ctx, docker.TaskRequest{
		Image: cfg.image,
		Cmd:   cfg.cmd,
		Input: input,
		Files: map[string][]byte{
			cfg.fileName: binary,
		},
		Limits: docker.Resources{
			Memory:  memoryLimit,
			Timeout: time.Duration(timeLimitNano),
		},
	})
	if err != nil {
		return "", err
	}

	if resp.ExitCode == 137 {
		return "", fmt.Errorf("resource limit exceeded (OOM or Time Limit)")
	}

	output := resp.Stdout
	if resp.Stderr != "" {
		output += "\nStderr:\n" + resp.Stderr
	}

	return output, nil
}
