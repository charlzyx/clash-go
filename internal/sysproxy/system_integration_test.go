//go:build integration

package sysproxy

import "testing"

// TestLiveSystemProxy is where a macOS run that really changes the system
// proxy belongs. The default test set does not build this file.
func TestLiveSystemProxy(t *testing.T) {
	t.Skip("M1 does not change the OS proxy")
}
