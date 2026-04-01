# ContestPlatform

Desktop contest platform built with Wails, Go and React.

## Build Commands

Use the root `Makefile`:

```bash
make help
make setup
make test
make dev
make build
```

Platform-specific production packaging:

```bash
make build-macos
make build-linux
make build-windows
```

Cross-platform tagged binaries for verification:

```bash
make cross-build
```

## Bundled Compilers

The application can load bundled compilers from a `compilers` directory located next to the packaged binary.
For custom paths, set:

```bash
COMPILERS_DIR=/absolute/path/to/compilers make build
```

## Important

- Run packaged Wails outputs from `build/bin`.
- Do not run raw binaries produced by plain `go build` without `-tags production`.
