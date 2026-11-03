package sysproxy

import "runtime"

// Platform is the controller the app uses when the caller does not pass one.
// macOS uses networksetup/scutil. Linux uses GNOME gsettings and refuses
// other desktops. Every other GOOS gets Unsupported, whose Set fails, so
// the switch cannot show a success the operating system did not apply.
// Tests pass their own Controller; Memory remains the in-process fake.
func Platform() Controller {
	switch runtime.GOOS {
	case "darwin":
		return NewSystem()
	case "linux":
		return NewGnome()
	default:
		return Unsupported{GOOS: runtime.GOOS}
	}
}
