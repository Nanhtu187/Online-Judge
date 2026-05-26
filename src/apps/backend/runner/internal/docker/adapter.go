package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type TaskRequest struct {
	Image       string
	Cmd         []string
	Input       string
	Files       map[string][]byte // Files to inject (name -> content)
	ExtractFile string            // Optional: single file to extract back
	Limits      Resources
}

type Resources struct {
	Memory  int64
	Timeout time.Duration
}

type TaskResponse struct {
	Stdout     string
	Stderr     string
	ExitCode   int64
	OutputFile []byte // Content of ExtractFile if requested
}

type Adapter interface {
	RunTask(ctx context.Context, req TaskRequest) (*TaskResponse, error)
	RemoveContainer(ctx context.Context, containerID string) error
}

type dockerAdapter struct {
	cli *client.Client
}

func NewAdapter() (Adapter, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerAdapter{cli: cli}, nil
}

func (a *dockerAdapter) RunTask(ctx context.Context, req TaskRequest) (*TaskResponse, error) {
	// 1. Create
	resp, err := a.cli.ContainerCreate(ctx, &container.Config{
		Image:        req.Image,
		Cmd:          req.Cmd,
		WorkingDir:   "/",
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory: req.Limits.Memory,
		},
	}, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("create fail: %w", err)
	}
	containerID := resp.ID
	defer func() {
		// Use a fresh context for cleanup to ensure it runs even if parent is cancelled
		ctx := context.Background()
		a.cli.ContainerKill(ctx, containerID, "SIGKILL")
		if err := a.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
			fmt.Printf("Warning: failed to remove container %s: %v\n", containerID, err)
		}
	}()

	// 2. Inject Files
	if len(req.Files) > 0 {
		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		for name, content := range req.Files {
			hdr := &tar.Header{
				Name: name,
				Mode: 0755,
				Size: int64(len(content)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return nil, err
			}
			if _, err := tw.Write(content); err != nil {
				return nil, err
			}
		}
		tw.Close()
		if err := a.cli.CopyToContainer(ctx, containerID, "/", &buf, container.CopyToContainerOptions{}); err != nil {
			return nil, fmt.Errorf("copy in fail: %w", err)
		}
	}

	// 3. Attach for Input
	attachResp, err := a.cli.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("attach fail: %w", err)
	}
	defer attachResp.Close()

	// 4. Start
	if err := a.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("start fail: %w", err)
	}

	// 5. Pipe Stdin
	if req.Input != "" {
		go func() {
			defer attachResp.CloseWrite()
			io.WriteString(attachResp.Conn, req.Input)
		}()
	} else {
		attachResp.CloseWrite()
	}

	// 6. Wait
	var exitCode int64
	waitCtx := ctx
	if req.Limits.Timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, req.Limits.Timeout)
		defer cancel()
	}

	statusCh, errCh := a.cli.ContainerWait(waitCtx, containerID, container.WaitConditionNotRunning)

	select {
	case <-waitCtx.Done():
		if waitCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("time limit exceeded")
		}
		return nil, waitCtx.Err()
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case status := <-statusCh:
		exitCode = status.StatusCode
	}

	// 7. Collect Logs
	logReader, err := a.cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, fmt.Errorf("logs fail: %w", err)
	}
	defer logReader.Close()

	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, logReader)

	res := &TaskResponse{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}

	// 8. Extract File if needed
	if req.ExtractFile != "" && exitCode == 0 {
		reader, _, err := a.cli.CopyFromContainer(ctx, containerID, "/"+req.ExtractFile)
		if err == nil {
			defer reader.Close()
			tr := tar.NewReader(reader)
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}
				if strings.HasSuffix(hdr.Name, req.ExtractFile) {
					res.OutputFile, _ = io.ReadAll(tr)
					break
				}
			}
		}
	}

	return res, nil
}

func (a *dockerAdapter) RemoveContainer(ctx context.Context, containerID string) error {
	return a.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}
