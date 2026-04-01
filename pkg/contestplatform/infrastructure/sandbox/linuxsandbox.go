//go:build linux || darwin

package sandbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
)

func NewLinuxSandbox() appmodel.Sandbox {
	return &linuxSandbox{}
}

type linuxSandbox struct{}

func (l *linuxSandbox) Prepare(ctx context.Context, lang domainmodel.Language, sourceCode string) (string, string, error) {
	cfg, ok := appmodel.Languages[lang]
	if !ok {
		return "", "", fmt.Errorf("language %s not supported", lang)
	}

	tmpDir, err := os.MkdirTemp("", "contestplatform-*")
	if err != nil {
		return "", "", err
	}

	sourcePath := filepath.Join(tmpDir, cfg.SourceFile)
	err = os.WriteFile(sourcePath, []byte(sourceCode), 0644)
	if err != nil {
		return "", "", err
	}

	if cfg.IsCompiled() {
		command, err := cfg.CompileCommand(tmpDir)
		if err != nil {
			return tmpDir, "", err
		}

		cmd := exec.CommandContext(ctx, command.Path, command.Args...)
		cmd.Dir = tmpDir
		cmd.Env = command.Env
		output, runErr := cmd.CombinedOutput()
		if runErr != nil {
			return tmpDir, "", fmt.Errorf("compilation error: %s", string(output))
		}

		runCommand, err := cfg.RunCommand(tmpDir)
		if err != nil {
			return tmpDir, "", err
		}

		return tmpDir, runCommand.Path, nil
	}

	runCommand, err := cfg.RunCommand(tmpDir)
	if err != nil {
		return tmpDir, "", err
	}

	launcherPath := filepath.Join(tmpDir, "run.sh")
	script := fmt.Sprintf("#!/bin/sh\nexec \"%s\"",
		strings.ReplaceAll(runCommand.Path, "\"", "\\\""),
	)
	for _, arg := range runCommand.Args {
		script += fmt.Sprintf(" \"%s\"", strings.ReplaceAll(arg, "\"", "\\\""))
	}
	script += "\n"

	if err = os.WriteFile(launcherPath, []byte(script), 0755); err != nil {
		return "", "", err
	}

	return tmpDir, launcherPath, nil
}

func (l *linuxSandbox) Execute(
	ctx context.Context,
	dir, exe string,
	input string,
	limits domainmodel.Constraints,
) (appmodel.SandboxResponse, error) {
	if limits.TimeLimit > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, limits.TimeLimit)
		defer cancel()
	}

	runCommand := exe
	runArgs := []string(nil)
	if !filepath.IsAbs(exe) {
		runCommand = filepath.Join(dir, exe)
	}
	if _, statErr := os.Stat(runCommand); statErr != nil {
		return appmodel.SandboxResponse{}, statErr
	}

	cmd := exec.CommandContext(ctx, runCommand, runArgs...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	cmd.Env = appmodel.ToolchainEnv(dir)

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
		} else if errors.Is(err, context.DeadlineExceeded) {
			exitCode = -1
		} else {
			return appmodel.SandboxResponse{}, err
		}
	}

	return appmodel.SandboxResponse{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		TimeUsed:  duration,
		ExitCode:  exitCode,
		IsTimeout: errors.Is(ctx.Err(), context.DeadlineExceeded),
	}, nil
}

func (l *linuxSandbox) Cleanup(workingDir string) error {
	return os.RemoveAll(workingDir)
}
