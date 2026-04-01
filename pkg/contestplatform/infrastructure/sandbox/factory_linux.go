//go:build linux

package sandbox

import appmodel "contest-platform/pkg/contestplatform/app/model"

func NewSandbox() (appmodel.Sandbox, error) {
	return NewLinuxSandbox(), nil
}
