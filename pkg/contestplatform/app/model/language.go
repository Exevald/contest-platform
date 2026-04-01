package model

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"

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

func (cfg LanguageConfig) IsCompiled() bool {
	return cfg.CompilerBinary != ""
}

func (cfg LanguageConfig) CompileCommand() (string, []string) {
	if !cfg.IsCompiled() {
		return "", nil
	}

	args := make([]string, 0, len(cfg.CompilerArgs)+2)
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

	return resolveBundledBinary(cfg.CompilerBinary), args
}

func (cfg LanguageConfig) RunCommand() (string, []string) {
	if cfg.IsCompiled() {
		return cfg.ExecutableName(), nil
	}

	return resolveBundledBinary(cfg.InterpreterBinary), append([]string{cfg.SourceFile}, cfg.RunArgs...)
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

func resolveBundledBinary(binary string) string {
	if binary == "" {
		return ""
	}

	root := bundledCompilersDir()
	if root != "" {
		candidates := []string{
			filepath.Join(root, binary),
			filepath.Join(root, filepath.Base(binary)),
		}
		if runtime.GOOS == "windows" && filepath.Ext(binary) == "" {
			candidates = append(candidates,
				filepath.Join(root, binary+".exe"),
				filepath.Join(root, filepath.Base(binary)+".exe"),
			)
		}
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	return binary
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
