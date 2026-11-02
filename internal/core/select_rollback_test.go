package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charlzyx/clash-go/internal/ui"

	mygoui "github.com/egoist/mygo/ui"
)

// TestUIRollsBackMissingNode clicks a node the kernel does not have.
// The hook is core.Select, read back with Selected (the in-process Now),
// not PUT /proxies. Do not call t.Parallel: mihomo is a process singleton.
func TestUIRollsBackMissingNode(t *testing.T) {
	dir := t.TempDir()
	mixed := freePort(t)
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
dns:
  enable: false
proxies:
  - {name: n1, type: direct}
  - {name: n2, type: direct}
proxy-groups:
  - {name: PROXY, type: select, proxies: [n1, n2]}
rules:
  - MATCH,PROXY
`, mixed)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	c := New()
	if err := c.Start(dir); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = c.Stop() })
	if err := c.Select("PROXY", "n1"); err != nil {
		t.Fatal(err)
	}

	m := ui.NewModel(ui.HomeState{Running: true, Mode: "rule"})
	m.Profiles = []ui.ProfileItem{{UID: "p1", Name: "sub", Type: "remote", URL: "https://example.com/sub"}}
	m.CurrentProfile = "p1"
	m.Prefs.ProxyColumns = 1
	m.Groups = []ui.ProxyGroup{{
		Name: "PROXY", Type: "select", Now: "n1", Open: true,
		Nodes: []ui.ProxyNode{
			{Name: "n1", Type: "direct", Delay: -1},
			{Name: "ghost", Type: "direct", Delay: -1},
		},
	}}
	m.Home.Node = ui.CurrentNode{Group: "PROXY", Name: "n1", Type: "direct", Delay: -1}
	m.Hooks.SelectProxy = c.Select
	m.Router.Replace("/proxies")

	tt := mygoui.NewTester(m.View, ui.WindowWidth, ui.WindowHeight)
	tt.SetScale(1)
	for i := 0; i < 5 && tt.FramePending(); i++ {
		tt.Advance(time.Second)
	}
	if err := tt.Click("node-PROXY-ghost"); err != nil {
		t.Fatal(err)
	}
	if m.Groups[0].Now != "n1" || m.Home.Node.Name != "n1" {
		t.Fatalf("UI kept the missing node: now %q home %q notice %q", m.Groups[0].Now, m.Home.Node.Name, m.Notice)
	}
	if m.Notice == "" {
		t.Fatal("missing node did not surface an error")
	}
	if !tt.HasText(m.Notice) {
		t.Fatalf("notice %q is not on the page, texts %q", m.Notice, tt.Texts())
	}
	got, err := c.Selected("PROXY")
	if err != nil {
		t.Fatal(err)
	}
	if got != "n1" {
		t.Fatalf("kernel Now %q, want n1", got)
	}
}
