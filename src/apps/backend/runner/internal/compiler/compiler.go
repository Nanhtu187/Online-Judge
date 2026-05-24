package compiler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/docker"
)

type Compiler interface {
	Compile(ctx context.Context, lang string, sourceCode string) ([]byte, error)
	IsCompileLanguage(lang string) bool
}

type DockerCompiler struct {
	adapter docker.Adapter
}

func NewDockerCompiler(adapter docker.Adapter) *DockerCompiler {
	return &DockerCompiler{adapter: adapter}
}

type config struct {
	image      string
	fileName   string
	cmd        []string
	outputName string
}

var langConfigs = map[string]config{
	"go": {
		image:      "golang:1.24-alpine",
		fileName:   "main.go",
		cmd:        []string{"go", "build", "-o", "output", "main.go"},
		outputName: "output",
	},
	"cpp": {
		image:      "gcc:latest",
		fileName:   "main.cpp",
		cmd:        []string{"g++", "-static", "-o", "output", "main.cpp"},
		outputName: "output",
	},
}

func (c *DockerCompiler) IsCompileLanguage(lang string) bool {
	_, ok := langConfigs[strings.ToLower(lang)]
	return ok
}

func (c *DockerCompiler) Compile(ctx context.Context, lang string, sourceCode string) ([]byte, error) {
	cfg, ok := langConfigs[strings.ToLower(lang)]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	resp, err := c.adapter.RunTask(ctx, docker.TaskRequest{
		Image: cfg.image,
		Cmd:   cfg.cmd,
		Files: map[string][]byte{
			cfg.fileName: []byte(sourceCode),
		},
		ExtractFile: cfg.outputName,
	})
	if err != nil {
		return nil, err
	}

	if resp.ExitCode != 0 {
		return nil, fmt.Errorf("compilation failed (exit %d):\nStdout: %s\nStderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	if len(resp.OutputFile) == 0 {
		return nil, fmt.Errorf("binary not found in compilation output")
	}

	return resp.OutputFile, nil
}
