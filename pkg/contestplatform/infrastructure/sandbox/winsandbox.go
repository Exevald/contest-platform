//go:build windows

package sandbox

import (
	"bytes"
	"context"
	stderror "errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
)

func NewWindowsSandbox() appmodel.Sandbox {
	return &winSandbox{}
}

type winSandbox struct {
}

func (w *winSandbox) Prepare(ctx context.Context, lang domainmodel.Language, sourceCode string) (string, string, error) {
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

	exePath := filepath.Join(tmpDir, cfg.ExecutableName())
	if cfg.IsCompiled() {
		command, args := cfg.CompileCommand()
		cmd := exec.CommandContext(ctx, command, args...)
		cmd.Dir = tmpDir
		out, err2 := cmd.CombinedOutput()
		if err2 != nil {
			return tmpDir, "", fmt.Errorf("compilation error: %s", string(out))
		}
		return tmpDir, exePath, nil
	}

	interpreter, args := cfg.RunCommand()
	launcherPath := filepath.Join(tmpDir, "run.cmd")
	commandLine := "@" + "echo off\r\n\"" + interpreter + "\""
	for _, arg := range args {
		commandLine += " \"" + arg + "\""
	}
	commandLine += "\r\n"

	if err = os.WriteFile(launcherPath, []byte(commandLine), 0644); err != nil {
		return "", "", err
	}

	return tmpDir, launcherPath, nil
}

func (w *winSandbox) Execute(ctx context.Context, dir, exe string, input string, limits domainmodel.Constraints) (appmodel.SandboxResponse, error) {
	if limits.TimeLimit > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, limits.TimeLimit)
		defer cancel()
	}

	cfgExe := exe
	if !filepath.IsAbs(cfgExe) {
		cfgExe = filepath.Join(dir, cfgExe)
	}
	cmd := exec.CommandContext(ctx, "cmd", "/C", cfgExe)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if stderror.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		} else if stderror.Is(err, context.DeadlineExceeded) {
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
		IsTimeout: stderror.Is(ctx.Err(), context.DeadlineExceeded),
	}, nil
}

func (w *winSandbox) Cleanup(workingDir string) error {
	return os.RemoveAll(workingDir)
}
