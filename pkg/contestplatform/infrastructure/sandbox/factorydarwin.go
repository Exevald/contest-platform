//go:build darwin

package sandbox

import appmodel "contest-platform/pkg/contestplatform/app/model"

func NewSandbox() (appmodel.Sandbox, error) {
	return NewDarwinSandbox(), nil
}
