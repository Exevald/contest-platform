package model

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveToolBinaryFromBundledCompilersDir(t *testing.T) {
	tempDir := t.TempDir()
	platformDir := filepath.Join(tempDir, runtime.GOOS+"-"+runtime.GOARCH, "bin")
	if err := os.MkdirAll(platformDir, 0755); err != nil {
		t.Fatalf("mkdir toolchain dir: %v", err)
	}

	binaryName := "python"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(platformDir, binaryName)
	if err := os.WriteFile(binaryPath, []byte(""), 0755); err != nil {
		t.Fatalf("write bundled tool: %v", err)
	}

	t.Setenv(BundledCompilersDirEnv, tempDir)

	resolved, err := resolveToolBinary("python")
	if err != nil {
		t.Fatalf("resolve tool: %v", err)
	}
	if resolved != binaryPath {
		t.Fatalf("expected %s, got %s", binaryPath, resolved)
	}
}

func TestToolchainEnvPrependsCompilersPath(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(BundledCompilersDirEnv, tempDir)

	env := ToolchainEnv(t.TempDir())

	var pathValue string
	for _, item := range env {
		if strings.HasPrefix(item, "PATH=") {
			pathValue = strings.TrimPrefix(item, "PATH=")
			break
		}
	}

	if pathValue == "" {
		t.Fatal("PATH was not set")
	}
	if !strings.Contains(pathValue, tempDir) {
		t.Fatalf("expected PATH to contain %s, got %s", tempDir, pathValue)
	}
}

func TestLanguageConfigIsAvailable(t *testing.T) {
	tempDir := t.TempDir()
	platformDir := filepath.Join(tempDir, runtime.GOOS+"-"+runtime.GOARCH, "bin")
	if err := os.MkdirAll(platformDir, 0755); err != nil {
		t.Fatalf("mkdir toolchain dir: %v", err)
	}

	binaryName := "custom-compiler"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(platformDir, binaryName)
	if err := os.WriteFile(binaryPath, []byte(""), 0755); err != nil {
		t.Fatalf("write toolchain binary: %v", err)
	}

	t.Setenv(BundledCompilersDirEnv, tempDir)

	if !(LanguageConfig{CompilerBinary: "custom-compiler"}).IsAvailable() {
		t.Fatal("expected compiler to be available")
	}
	if (LanguageConfig{CompilerBinary: "missing-compiler"}).IsAvailable() {
		t.Fatal("expected missing compiler to be unavailable")
	}
}

func TestSystemBinaryCandidatesIncludesPython3Alias(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("python3 alias is unix-specific")
	}

	candidates := systemBinaryCandidates("python")
	joined := strings.Join(candidates, "\n")
	if !strings.Contains(joined, "python3") {
		t.Fatalf("expected python3 fallback in candidates, got %v", candidates)
	}
}
