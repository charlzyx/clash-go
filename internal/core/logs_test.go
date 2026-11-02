package core

import (
	"strings"
	"testing"
	"time"

	"github.com/metacubex/mihomo/log"
)

func TestLogsSinceRequiresRunningCore(t *testing.T) {
	c := New()
	if _, _, ok := c.LogsSince(0); ok {
		t.Fatal("stopped core reported logs")
	}
	if c.ClearLogs() {
		t.Fatal("stopped core cleared logs")
	}
}

func TestLogFeedRecordsAndClears(t *testing.T) {
	dir := t.TempDir()
	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = c.Stop() })

	const payload = "clash-go-log-feed-marker"
	log.Infoln("%s", payload)

	deadline := time.Now().Add(2 * time.Second)
	var seq uint64
	found := false
	for time.Now().Before(deadline) {
		lines, next, ok := c.LogsSince(0)
		if !ok {
			t.Fatal("running core did not report logs")
		}
		seq = next
		for _, line := range lines {
			if line.Level == "info" && strings.Contains(line.Payload, payload) && !line.Time.IsZero() {
				found = true
			}
		}
		if found {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !found {
		t.Fatal("marker log was not recorded")
	}

	if !c.ClearLogs() {
		t.Fatal("clear on a running core returned false")
	}
	lines, next, ok := c.LogsSince(seq)
	if !ok {
		t.Fatal("logs unavailable after clear")
	}
	for _, line := range lines {
		if strings.Contains(line.Payload, payload) {
			t.Fatal("clear left the marker")
		}
	}
	if next < seq {
		t.Fatalf("cursor went backwards: %d -> %d", seq, next)
	}

	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := c.LogsSince(0); ok {
		t.Fatal("stopped core still reports logs")
	}
}
