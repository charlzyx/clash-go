package sysproxy

import (
	"fmt"
	"runtime"
)

// Unsupported is the controller for systems without a real system-proxy
// backend (Windows in this tree). Set fails so the UI can roll the switch
// back. Remembering "on" without touching the OS is a false success.
type Unsupported struct {
	GOOS string
}

// Current reports the proxy as off.
func (u Unsupported) Current() (Settings, error) { return Settings{}, nil }

// Set refuses the change.
func (u Unsupported) Set(Settings) error {
	name := u.GOOS
	if name == "" {
		name = runtime.GOOS
	}
	return fmt.Errorf("当前系统（%s）尚未支持系统代理", name)
}
