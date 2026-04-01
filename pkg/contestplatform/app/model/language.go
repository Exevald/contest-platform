package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	domainmodel "contest-platform/pkg/contestplatform/domain/model"
)

const BundledCompilersDirEnv = "CONTESTPLATFORM_COMPILERS_DIR"

type LanguageConfig struct {
	DisplayName       string
	SourceFile        string
	SourceExt         string
	CompilerBinary    string
	CompilerArgs      []string
	InterpreterBinary string
	RunArgs           []string
	ExeFile           string
}

type ToolchainCommand struct {
	Path string
	Args []string
	Env  []string
}

func ToolchainEnv(workDir string) []string {
	return toolchainEnv(workDir)
}

func (cfg LanguageConfig) IsCompiled() bool {
	return cfg.CompilerBinary != ""
}

func (cfg LanguageConfig) CompileCommand(workDir string) (ToolchainCommand, error) {
	if !cfg.IsCompiled() {
		return ToolchainCommand{}, nil
	}

	commandPath, err := resolveToolBinary(cfg.CompilerBinary)
	if err != nil {
		return ToolchainCommand{}, err
	}

	args := make([]string, 0, len(cfg.CompilerArgs))
	for _, arg := range cfg.CompilerArgs {
		switch arg {
		case "{source}":
			args = append(args, cfg.SourceFile)
		case "{exe}":
			args = append(args, cfg.ExecutableName())
		default:
			args = append(args, arg)
		}
	}

	return ToolchainCommand{
		Path: commandPath,
		Args: args,
		Env:  toolchainEnv(workDir),
	}, nil
}

func (cfg LanguageConfig) RunCommand(workDir string) (ToolchainCommand, error) {
	if cfg.IsCompiled() {
		return ToolchainCommand{
			Path: filepath.Join(workDir, cfg.ExecutableName()),
			Env:  toolchainEnv(workDir),
		}, nil
	}

	commandPath, err := resolveToolBinary(cfg.InterpreterBinary)
	if err != nil {
		return ToolchainCommand{}, err
	}

	return ToolchainCommand{
		Path: commandPath,
		Args: append([]string{cfg.SourceFile}, cfg.RunArgs...),
		Env:  toolchainEnv(workDir),
	}, nil
}

func (cfg LanguageConfig) ExecutableName() string {
	if cfg.ExeFile == "" {
		return ""
	}
	if runtime.GOOS == "windows" && filepath.Ext(cfg.ExeFile) == "" {
		return cfg.ExeFile + ".exe"
	}
	return cfg.ExeFile
}

func resolveToolBinary(binary string) (string, error) {
	if binary == "" {
		return "", nil
	}

	for _, candidate := range binaryCandidates(binary) {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	resolved, err := exec.LookPath(binary)
	if err == nil {
		return resolved, nil
	}

	for _, candidate := range binaryCandidates(binary) {
		if resolved, err = exec.LookPath(candidate); err == nil {
			return resolved, nil
		}
	}

	return "", fmt.Errorf("unable to resolve tool %q in compilers directory or PATH", binary)
}

func binaryCandidates(binary string) []string {
	binaries := []string{binary, filepath.Base(binary)}
	if runtime.GOOS == "windows" && filepath.Ext(binary) == "" {
		binaries = append(binaries, binary+".exe", filepath.Base(binary)+".exe")
	}

	seen := make(map[string]struct{}, len(binaries)*len(toolchainRoots()))
	candidates := make([]string, 0, len(binaries)*len(toolchainRoots()))
	for _, root := range toolchainRoots() {
		for _, name := range binaries {
			candidate := filepath.Join(root, name)
			if _, ok := seen[candidate]; ok {
				continue
			}
			seen[candidate] = struct{}{}
			candidates = append(candidates, candidate)
		}
	}

	return candidates
}

func toolchainRoots() []string {
	root := bundledCompilersDir()
	if root == "" {
		return nil
	}

	entries := []string{
		root,
		filepath.Join(root, "bin"),
		filepath.Join(root, runtime.GOOS),
		filepath.Join(root, runtime.GOOS, "bin"),
		filepath.Join(root, runtime.GOOS+"-"+runtime.GOARCH),
		filepath.Join(root, runtime.GOOS+"-"+runtime.GOARCH, "bin"),
	}

	seen := map[string]struct{}{}
	roots := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry == "" {
			continue
		}
		if _, ok := seen[entry]; ok {
			continue
		}
		seen[entry] = struct{}{}
		roots = append(roots, entry)
	}

	return roots
}

func toolchainPath() string {
	parts := make([]string, 0, len(toolchainRoots()))
	for _, root := range toolchainRoots() {
		if info, err := os.Stat(root); err == nil && info.IsDir() {
			parts = append(parts, root)
		}
	}

	return strings.Join(parts, string(os.PathListSeparator))
}

func toolchainEnv(workDir string) []string {
	env := os.Environ()
	extraPath := toolchainPath()
	if extraPath != "" {
		env = upsertEnv(env, "PATH", extraPath+string(os.PathListSeparator)+os.Getenv("PATH"))
	}

	cacheDir := filepath.Join(workDir, ".toolcache")
	goCacheDir := filepath.Join(cacheDir, "go-build")
	goModCacheDir := filepath.Join(cacheDir, "go-mod")
	goTmpDir := filepath.Join(cacheDir, "tmp")
	_ = os.MkdirAll(goCacheDir, 0755)
	_ = os.MkdirAll(goModCacheDir, 0755)
	_ = os.MkdirAll(goTmpDir, 0755)

	env = upsertEnv(env, "GOCACHE", goCacheDir)
	env = upsertEnv(env, "GOMODCACHE", goModCacheDir)
	env = upsertEnv(env, "GOTMPDIR", goTmpDir)
	env = upsertEnv(env, "GO111MODULE", "off")
	env = upsertEnv(env, "CGO_ENABLED", "0")

	return env
}

func upsertEnv(env []string, key string, value string) []string {
	prefix := key + "="
	for index, item := range env {
		if strings.HasPrefix(item, prefix) {
			env[index] = prefix + value
			return env
		}
	}

	return append(env, prefix+value)
}

func bundledCompilersDir() string {
	if custom := os.Getenv(BundledCompilersDirEnv); custom != "" {
		return custom
	}

	executable, err := os.Executable()
	if err != nil {
		return filepath.Join(".", "compilers")
	}

	return filepath.Join(filepath.Dir(executable), "compilers")
}

var Languages = map[domainmodel.Language]LanguageConfig{
	"cpp": {
		DisplayName:    "C++",
		SourceFile:     "solution.cpp",
		SourceExt:      "cpp",
		CompilerBinary: "g++",
		CompilerArgs:   []string{"{source}", "-O2", "-o", "{exe}"},
		ExeFile:        "solution",
	},
	"go": {
		DisplayName:    "Go",
		SourceFile:     "main.go",
		SourceExt:      "go",
		CompilerBinary: "go",
		CompilerArgs:   []string{"build", "-o", "{exe}", "{source}"},
		ExeFile:        "solution",
	},
	"js": {
		DisplayName:       "JavaScript",
		SourceFile:        "solution.js",
		SourceExt:         "js",
		InterpreterBinary: "node",
	},
	"pascal": {
		DisplayName:    "Pascal",
		SourceFile:     "solution.pas",
		SourceExt:      "pas",
		CompilerBinary: "fpc",
		CompilerArgs:   []string{"{source}", "-o{exe}"},
		ExeFile:        "solution",
	},
	"php": {
		DisplayName:       "PHP",
		SourceFile:        "solution.php",
		SourceExt:         "php",
		InterpreterBinary: "php",
	},
	"python": {
		DisplayName:       "Python",
		SourceFile:        "solution.py",
		SourceExt:         "py",
		InterpreterBinary: "python",
	},
}

type UILanguage struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
}

func SupportedUILanguages() []UILanguage {
	languages := make([]UILanguage, 0, len(Languages))
	for _, config := range Languages {
		languages = append(languages, UILanguage{
			Name:      config.DisplayName,
			Extension: config.SourceExt,
		})
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Name < languages[j].Name
	})

	return languages
}
