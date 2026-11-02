package sitetest

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProbeURL measures one GET. client nil uses a 5 second timeout.
// Callers in tests pass an httptest transport so this does not dial the public network.
func ProbeURL(raw string, client *http.Client) (int, error) {
	if raw == "" {
		return 0, fmt.Errorf("empty url")
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	req, err := http.NewRequest(http.MethodGet, raw, nil)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}
	ms := int(time.Since(start).Milliseconds())
	if ms < 0 {
		ms = 0
	}
	return ms, nil
}
