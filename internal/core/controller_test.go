//go:build !race

package core

import (
	"testing"
)

func TestExternalControllerListens(t *testing.T) {
	dir := t.TempDir()
	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Stop() })
	if err := waitListen("127.0.0.1:9097"); err != nil {
		t.Fatal(err)
	}
}
