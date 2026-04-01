package model

import (
	"challenger/pkg/challenger/domain/model"
)

type LanguageConfig struct {
	IsCompiled  bool
	SourceFile  string
	CompilerCmd string
	RunCmd      []string
	ExeFile     string
	SourceExt   string
}

var Languages = map[model.Language]LanguageConfig{
	"cpp": {
		IsCompiled:  true,
		SourceFile:  "solution.cpp",
		CompilerCmd: "g++ solution.cpp -O2 -o solution.out",
		RunCmd:      []string{"./solution.out"},
		ExeFile:     "solution.out",
	},
	"python": {
		IsCompiled: false,
		SourceFile: "solution.py",
		RunCmd:     []string{"python", "solution.py"},
	},
	"pascal": {
		IsCompiled:  true,
		SourceFile:  "solution.pas",
		CompilerCmd: "fpc solution.pas",
		ExeFile:     "solution",
		RunCmd:      []string{"./solution"},
	},
	"js": {
		IsCompiled: false,
		SourceFile: "solution.js",
		RunCmd:     []string{"node", "solution.js"},
	},
	"php": {
		IsCompiled: false,
		SourceFile: "solution.php",
		RunCmd:     []string{"php", "solution.php"},
	},
	"go": {
		IsCompiled:  true,
		SourceFile:  "main.go",
		CompilerCmd: "go build -o solution.out main.go",
		ExeFile:     "solution.out",
		RunCmd:      []string{"./solution.out"},
	},
}
