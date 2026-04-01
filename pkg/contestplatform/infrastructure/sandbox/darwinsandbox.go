//go:build darwin

package sandbox

import appmodel "contest-platform/pkg/contestplatform/app/model"

func NewDarwinSandbox() appmodel.Sandbox {
	return &linuxSandbox{}
}
