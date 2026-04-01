//go:build !linux && !darwin && !windows

package sandbox

import (
	"errors"

	appmodel "contest-platform/pkg/contestplatform/app/model"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")

func NewSandbox() (appmodel.Sandbox, error) {
	return nil, ErrUnsupportedPlatform
}
