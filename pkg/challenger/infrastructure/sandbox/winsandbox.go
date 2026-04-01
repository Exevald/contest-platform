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

	appmodel "challenger/pkg/challenger/app/model"
	"challenger/pkg/challenger/domain/model"
)

func NewWindowsSandbox() appmodel.Sandbox {
	return &winSandbox{}
}

type winSandbox struct {
}

func (w *winSandbox) Prepare(ctx context.Context, lang model.Language, sourceCode string) (string, string, error) {
	cfg, ok := appmodel.Languages[lang]
	if !ok {
		return "", "", fmt.Errorf("language %s not supported", lang)
	}

	tmpDir, err := os.MkdirTemp("", "challenger-*")
	if err != nil {
		return "", "", err
	}

	sourcePath := filepath.Join(tmpDir, cfg.SourceFile)
	err = os.WriteFile(sourcePath, []byte(sourceCode), 0644)
	if err != nil {
		return "", "", err
	}

	if cfg.IsCompiled {
		parts := strings.Split(cfg.CompilerCmd, " ")
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		cmd.Dir = tmpDir
		out, err2 := cmd.CombinedOutput()
		if err2 != nil {
			return tmpDir, "", fmt.Errorf("compilation error: %s", string(out))
		}
		return tmpDir, filepath.Join(tmpDir, cfg.ExeFile), nil
	}

	return tmpDir, "", nil
}

func (w *winSandbox) Execute(ctx context.Context, dir, exe string, input string, limits model.Constraints) (appmodel.SandboxResponse, error) {
	cmd := exec.CommandContext(ctx, "cmd", "/C", exe)
	if exe == "" {
	}

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
