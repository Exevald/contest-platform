package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	appmodel "challenger/pkg/challenger/app/model"
	"challenger/pkg/challenger/domain/model"
)

func NewLinuxSandbox() appmodel.Sandbox {
	return &linuxSandbox{}
}

type linuxSandbox struct{}

func (l *linuxSandbox) Prepare(ctx context.Context, lang model.Language, sourceCode string) (string, string, error) {
	cfg := appmodel.Languages[lang]
	tmpDir, _ := os.MkdirTemp("", "challenger-*")

	sourcePath := filepath.Join(tmpDir, cfg.SourceFile)
	err := os.WriteFile(sourcePath, []byte(sourceCode), 0644)
	if err != nil {
		return "", "", err
	}

	if cfg.IsCompiled {
		parts := strings.Split(cfg.CompilerCmd, " ")
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		cmd.Dir = tmpDir
		err = cmd.Run()
		if err != nil {
			return "", "", err
		}
		return tmpDir, filepath.Join(tmpDir, cfg.ExeFile), nil
	}
	return tmpDir, "", nil
}

func (l *linuxSandbox) Execute(
	ctx context.Context,
	dir, exe string,
	input string,
	limits model.Constraints,
) (appmodel.SandboxResponse, error) {
	cmd := exec.CommandContext(ctx, exe)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}

	return appmodel.SandboxResponse{
		Stdout:    stdout.String(),
		TimeUsed:  duration,
		ExitCode:  exitCode,
		IsTimeout: errors.Is(ctx.Err(), context.DeadlineExceeded),
	}, nil
}

func (l *linuxSandbox) Cleanup(workingDir string) error {
	return os.RemoveAll(workingDir)
}
