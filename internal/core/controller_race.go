//go:build race

package core

// startExternalController is a no-op when the race detector is on.
// mihomo v1.19.33 publishes the external controller from a goroutine and
// shuts it down by writing the same globals, which the detector reports.
func startExternalController(addr, secret string, debug bool) {}
