//go:build integration

package tun

import "testing"

// TestLiveTUN is where a run that really opens a TUN device belongs.
func TestLiveTUN(t *testing.T) {
	t.Skip("M1 does not open a TUN device")
}
