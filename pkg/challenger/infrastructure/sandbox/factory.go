package sandbox

import (
	stderr "errors"
	"runtime"

	"github.com/pkg/errors"

	appmodel "challenger/pkg/challenger/app/model"
)

var (
	ErrUnsupportedPlatform = stderr.New("unsupported platform")
)

func NewSandbox() (appmodel.Sandbox, error) {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsSandbox(), nil
	case "linux":
		return NewLinuxSandbox(), nil
	default:
		return &linuxSandbox{}, errors.WithStack(ErrUnsupportedPlatform)
	}
}
