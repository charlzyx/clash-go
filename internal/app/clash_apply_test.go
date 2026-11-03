package app

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/profile"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
	"github.com/charlzyx/clash-go/internal/ui"
)

// TestSavePrefsReloadsKernel proves a settings change is patched into
// config.yaml and applied by the running mihomo, not only stored as JSON.
// Do not call t.Parallel: the kernel is a process-wide singleton.
func TestSavePrefsReloadsKernel(t *testing.T) {
	dir := t.TempDir()
	mixed := freePort(t)
	mixed2 := freePort(t)
	socks := freePort(t)
	httpPort := freePort(t)
	ext := freePort(t)
	raw := []byte(strings.ReplaceAll(`mixed-port: MIXED
socks-port: 0
port: 0
allow-lan: false
mode: rule
log-level: info
ipv6: true
unified-delay: true
find-process-mode: off
geo-auto-update: false
external-controller: 127.0.0.1:EXT
dns:
  enable: false
proxies: []
rules:
  - MATCH,DIRECT
`, "MIXED", itoa(mixed)))
	raw = []byte(strings.ReplaceAll(string(raw), "EXT", itoa(ext)))
	profiles := profile.NewStore(dir)
	item, err := profiles.ImportLocal("sub.yaml", raw)
	if err != nil {
		t.Fatal(err)
	}
	prefs := ui.DefaultPrefs()
	prefs.MixedPort = mixed
	prefs.SocksOn = false
	prefs.HTTPOn = false
	prefs.RedirOn = false
	prefs.TProxyOn = false
	blob, err := json.Marshal(prefs)
	if err != nil {
		t.Fatal(err)
	}
	store := appconfig.NewMemory(dir)
	if err := store.Save(appconfig.Settings{Prefs: string(blob)}); err != nil {
		t.Fatal(err)
	}
	kernel := core.New()
	a := New(Options{
		Dirs:  appdir.Fixed(dir),
		Store: store,
		Proxy: sysproxy.NewMemory(),
		TUN:   tun.NewMemory(),
		Core:  kernel,
	})
	t.Cleanup(func() { _ = kernel.Stop() })
	if err := a.Prepare(); err != nil {
		t.Fatal(err)
	}
	if !kernel.Status().Running {
		t.Fatal(a.Model().Home.Error)
	}
	if a.Model().CurrentProfile != item.UID {
		t.Fatalf("current %q, imported %q", a.Model().CurrentProfile, item.UID)
	}

	next := a.Model().Prefs
	next.AllowLan = true
	next.IPv6 = false
	next.UnifiedDelay = false
	next.LogLevel = "warn"
	next.MixedPort = mixed2
	next.SocksOn = true
	next.SocksPort = socks
	next.HTTPOn = true
	next.HTTPPort = httpPort
	next.RedirOn = false
	next.TProxyOn = false
	next.DNSOverwrite = true
	if err := a.savePrefs(next); err != nil {
		t.Fatal(err)
	}

	saved, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(saved.Prefs, `"AllowLan":true`) {
		t.Fatalf("settings json missing allow-lan: %s", saved.Prefs)
	}
	body, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, want := range []string{
		"allow-lan: true",
		"ipv6: false",
		"unified-delay: false",
		"log-level: warning",
		"mixed-port: " + itoa(mixed2),
		"socks-port: " + itoa(socks),
		"port: " + itoa(httpPort),
		"enable: true",
		"8.8.8.8",
		"enhanced-mode: fake-ip",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("config missing %q\n%s", want, text)
		}
	}
	if runtime.GOOS != "windows" && !strings.Contains(text, "redir-port: 0") {
		t.Errorf("redir was not cleared\n%s", text)
	}
	if strings.Contains(text, "log-level: warn\n") {
		t.Fatalf("kernel would reject log-level warn\n%s", text)
	}
	profileBody, err := os.ReadFile(filepath.Join(dir, "profiles", item.UID+".yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(profileBody), "8.8.8.8") || strings.Contains(string(profileBody), "allow-lan: true") {
		t.Fatalf("subscription file was rewritten:\n%s", profileBody)
	}
	catalog, err := profiles.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Items) != 1 || !catalog.Items[0].Option.DNSOverwrite {
		t.Fatalf("dns flag not stored on the profile: %+v", catalog.Items)
	}

	snap := kernel.GeneralSnapshot()
	if !snap.Running || !snap.AllowLan || snap.IPv6 || snap.UnifiedDelay {
		t.Fatalf("kernel general %+v", snap)
	}
	if snap.LogLevel != "warning" {
		t.Fatalf("kernel log level %q", snap.LogLevel)
	}
	if snap.MixedPort != mixed2 || snap.SocksPort != socks || snap.Port != httpPort {
		t.Fatalf("kernel ports mixed=%d socks=%d http=%d", snap.MixedPort, snap.SocksPort, snap.Port)
	}
	if !snap.DNS {
		t.Fatal("kernel DNS resolver is off after overwrite")
	}

	cfgPath := filepath.Join(dir, "config.yaml")
	stamped, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, append(stamped, []byte("# sentinel\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	theme := a.Model().Prefs
	theme.ThemeMode = "dark"
	if err := a.savePrefs(theme); err != nil {
		t.Fatal(err)
	}
	afterTheme, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(afterTheme), "# sentinel") {
		t.Fatal("saving the theme rewrote the runtime config")
	}

	off := a.Model().Prefs
	off.DNSOverwrite = false
	if err := a.savePrefs(off); err != nil {
		t.Fatal(err)
	}
	offSnap := kernel.GeneralSnapshot()
	if offSnap.DNS || !offSnap.AllowLan || offSnap.LogLevel != "warning" {
		t.Fatalf("after dns off %+v", offSnap)
	}
	runtimeText, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(runtimeText), "8.8.8.8") {
		t.Fatalf("dns overwrite stayed in the runtime config:\n%s", runtimeText)
	}
	if !strings.Contains(string(runtimeText), "enable: false") {
		t.Fatalf("profile dns was not restored:\n%s", runtimeText)
	}
	catalog, err = profiles.Load()
	if err != nil {
		t.Fatal(err)
	}
	if catalog.Items[0].Option.DNSOverwrite {
		t.Fatal("profile still records dns overwrite")
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

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
