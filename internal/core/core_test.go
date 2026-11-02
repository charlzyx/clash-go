package core

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// These tests drive mihomo's process-wide singleton. Do not call t.Parallel.

func TestStartStopDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = c.Stop() })

	data, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{
		"mixed-port: 7897",
		"socks-port: 7898",
		"port: 7899",
		"external-controller: 127.0.0.1:9097",
		"MATCH,DIRECT",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("config missing %q\n%s", want, text)
		}
	}

	st := c.Status()
	if !st.Running {
		t.Fatal("status running = false")
	}
	if st.Version != EmbeddedVersion {
		t.Fatalf("version %q, want %s", st.Version, EmbeddedVersion)
	}
	if st.MixedPort != 7897 {
		t.Fatalf("mixed port %d, want 7897", st.MixedPort)
	}
	if st.ConfigDir != dir {
		t.Fatalf("config dir %q, want %q", st.ConfigDir, dir)
	}
	if st.Mode != "rule" {
		t.Fatalf("mode %q, want rule", st.Mode)
	}

	if err := waitListen("127.0.0.1:7897"); err != nil {
		t.Fatal(err)
	}
	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}
	if waitListen("127.0.0.1:7897") == nil {
		t.Fatal("mixed port still listening after Stop")
	}
	if c.Status().Running {
		t.Fatal("status still running")
	}
}

func TestStartTwice(t *testing.T) {
	dir := t.TempDir()
	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Stop() })
	if err := c.Start(dir); err == nil {
		t.Fatal("second Start succeeded")
	}
}

func TestHTTPSubscriptionDirect(t *testing.T) {
	const body = "proxies:\n  - name: direct-only\n    type: direct\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		_, _ = ioWrite(w, body)
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	mixed := freePort(t)
	api := freePort(t)
	cfg := fmt.Sprintf(`mixed-port: %d
socks-port: 0
port: 0
redir-port: 0
tproxy-port: 0
allow-lan: false
mode: rule
log-level: error
ipv6: false
find-process-mode: off
geo-auto-update: false
external-controller: 127.0.0.1:%d
dns:
  enable: false
proxy-providers:
  sub:
    type: http
    url: %q
    path: providers/sub.yaml
    interval: 86400
    health-check:
      enable: false
proxy-groups:
  - name: PROXY
    type: select
    use:
      - sub
    proxies:
      - DIRECT
rules:
  - MATCH,DIRECT
`, mixed, api, srv.URL)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = c.Stop() })

	nodes := c.ProviderProxies("sub")
	if len(nodes) != 1 {
		t.Fatalf("provider proxies %+v, want one DIRECT node", nodes)
	}
	if nodes[0].Name != "direct-only" || nodes[0].Type != "Direct" {
		t.Fatalf("node %+v, want direct-only Direct", nodes[0])
	}
	if err := c.Reload(); err != nil {
		t.Fatal(err)
	}
	nodes = c.ProviderProxies("sub")
	if len(nodes) != 1 || nodes[0].Name != "direct-only" {
		t.Fatalf("after reload %+v", nodes)
	}
	if c.Status().MixedPort != mixed {
		t.Fatalf("mixed port %d, want %d", c.Status().MixedPort, mixed)
	}
	if c.Status().ConfigDir != dir {
		t.Fatalf("config dir %q", c.Status().ConfigDir)
	}
}

func TestReloadWhenStopped(t *testing.T) {
	if err := New().Reload(); err == nil {
		t.Fatal("Reload on a stopped core succeeded")
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func waitListen(addr string) error {
	deadline := time.Now().Add(3 * time.Second)
	var last error
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		last = err
		time.Sleep(30 * time.Millisecond)
	}
	return fmt.Errorf("dial %s: %w", addr, last)
}

func ioWrite(w http.ResponseWriter, s string) (int, error) {
	return w.Write([]byte(s))
}

func TestConnectionSnapshotWhileStopped(t *testing.T) {
	c := New()
	up, down, items, ok := c.ConnectionSnapshot()
	if ok || up != 0 || down != 0 || items != nil {
		t.Fatalf("stopped snapshot ok=%v up=%d down=%d items=%v", ok, up, down, items)
	}
	if err := c.CloseAllConnections(); err == nil {
		t.Fatal("CloseAllConnections on a stopped core returned nil")
	}
	if err := c.CloseConnection("missing"); err == nil {
		t.Fatal("CloseConnection on a stopped core returned nil")
	}
}
