//go:build !race

package core

import "github.com/metacubex/mihomo/hub/route"

// externalControllerEnabled is false in the race-detector build, where
// startExternalController is a no-op. See controller_race.go.
const externalControllerEnabled = true

// startExternalController listens on the profile's external-controller
// after hub.Parse has set the log level. mihomo starts that server on a
// goroutine with no happens-before edge to Shutdown, so the race detector
// build leaves it down; see controller_race.go.
func startExternalController(addr, secret string, debug bool) {
	if addr == "" {
		return
	}
	route.ReCreateServer(&route.Config{
		Addr:    addr,
		Secret:  secret,
		IsDebug: debug,
	})
}
