# Compilers Layout

This directory is reserved for bundled compilers and interpreters used by the sandbox.

The runtime searches tools in this order:

1. `CONTESTPLATFORM_COMPILERS_DIR`
2. `./compilers`
3. `./compilers/bin`
4. `./compilers/<os>`
5. `./compilers/<os>/bin`
6. `./compilers/<os>-<arch>`
7. `./compilers/<os>-<arch>/bin`

Examples:

- `compilers/darwin-arm64/bin/python`
- `compilers/darwin-arm64/bin/node`
- `compilers/linux-amd64/bin/g++`
- `compilers/windows-amd64/bin/python.exe`

Supported tool names:

- `g++`
- `go`
- `node`
- `python`
- `php`
- `fpc`

The sandbox still falls back to the system `PATH` if a bundled tool is absent.
