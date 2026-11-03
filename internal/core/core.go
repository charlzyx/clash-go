package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/metacubex/mihomo/component/resolver"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	mihomodns "github.com/metacubex/mihomo/dns"
	"github.com/metacubex/mihomo/hub"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/listener"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel"
)

const configFileName = "config.yaml"

// Status is a snapshot of the running kernel. The UI reads this value
// and does not touch the mihomo singleton itself.
type Status struct {
	Running   bool
	Version   string
	MixedPort int
	Mode      string
	ConfigDir string
}

// ProxyNode is one proxy inside a provider.
type ProxyNode struct {
	Name string
	Type string
}

// Core controls the process-wide mihomo singleton.
type Core struct {
	mu      sync.Mutex
	running bool
	dir     string

	// logSub is the mihomo log subscription while the core is running.
	// logBook is the ring the UI drains. Both exist so logging never
	// blocks on a full subscriber: mihomo's Emit waits on that buffer.
	logSub  <-chan log.Event
	logDone chan struct{}
	logBook logBook
}

// New returns a core that is not running.
func New() *Core { return &Core{} }

// Start loads config.yaml from homeDir (writing DefaultConfig when the
// file is absent) and applies it with hub.Parse. homeDir becomes
// mihomo's configuration directory for this process.
func (c *Core) Start(homeDir string) error {
	if c == nil {
		return errors.New("core is nil")
	}
	if homeDir == "" {
		return errors.New("config dir is empty")
	}
	homeDir = filepath.Clean(homeDir)

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return errors.New("core already running")
	}
	C.SetHomeDir(homeDir)
	C.SetConfig(filepath.Join(homeDir, configFileName))

	// Subscribe before apply so startup lines are kept. Stop the feed if
	// apply fails; otherwise Emit can block the kernel logger on our buffer.
	c.startLogFeed()
	if err := apply(); err != nil {
		c.stopLogFeed()
		releaseListeners()
		return err
	}
	c.running = true
	c.dir = homeDir
	return nil
}

// Reload reads config.yaml from the directory passed to Start and applies
// it again.
func (c *Core) Reload() error {
	if c == nil {
		return errors.New("core is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return errors.New("core is not running")
	}
	return apply()
}

// Stop closes listeners, the external controller and provider update
// loops. It is safe to call more than once.
func (c *Core) Stop() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return nil
	}
	releaseOutboundProxies()
	releaseListeners()
	c.stopLogFeed()
	c.running = false
	return nil
}

// GeneralSnapshot is the clash fields the settings page writes, read back
// from the running kernel. DNS is true when mihomo installed a resolver,
// which is what dns.enable does. A stopped core returns a zero snapshot.
type GeneralSnapshot struct {
	Running      bool
	AllowLan     bool
	IPv6         bool
	UnifiedDelay bool
	LogLevel     string
	MixedPort    int
	SocksPort    int
	Port         int
	RedirPort    int
	TProxyPort   int
	DNS          bool
}

// GeneralSnapshot reads the kernel. It does not parse config.yaml, so a
// test can tell a reload from a JSON round-trip.
func (c *Core) GeneralSnapshot() GeneralSnapshot {
	var snap GeneralSnapshot
	if c == nil {
		return snap
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return snap
	}
	g := executor.GetGeneral()
	if g == nil {
		snap.Running = true
		return snap
	}
	snap.Running = true
	snap.AllowLan = g.AllowLan
	snap.IPv6 = g.IPv6
	snap.UnifiedDelay = g.UnifiedDelay
	snap.LogLevel = g.LogLevel.String()
	snap.MixedPort = g.MixedPort
	snap.SocksPort = g.SocksPort
	snap.Port = g.Port
	snap.RedirPort = g.RedirPort
	snap.TProxyPort = g.TProxyPort
	snap.DNS = resolver.DefaultResolver != nil
	return snap
}

// Status reports the kernel. Version is set even when the core is stopped.
func (c *Core) Status() Status {
	st := Status{Version: Version()}
	if c == nil {
		return st
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	st.Running = c.running
	st.ConfigDir = c.dir
	if !c.running {
		return st
	}
	general := executor.GetGeneral()
	st.MixedPort = general.MixedPort
	st.Mode = general.Mode.String()
	return st
}

// ProviderProxies returns the proxies loaded for a named provider.
// The slice is empty when the core is stopped or the provider is missing.
func (c *Core) ProviderProxies(name string) []ProxyNode {
	if c == nil || name == "" {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return nil
	}
	p, ok := tunnel.Providers()[name]
	if !ok || p == nil {
		return nil
	}
	proxies := p.Proxies()
	out := make([]ProxyNode, 0, len(proxies))
	for _, px := range proxies {
		out = append(out, ProxyNode{Name: px.Name(), Type: px.Type().String()})
	}
	return out
}

func apply() error {
	homeDir := C.Path.HomeDir()
	path := C.Path.Config()
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	if err := writeDefaultConfig(path); err != nil {
		return err
	}
	if err := config.Init(homeDir); err != nil {
		return fmt.Errorf("init config dir: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	// Parse once for the controller address, then hub.Parse a copy that
	// does not start it. mihomo v1.19.33's hub.ApplyConfig launches the
	// external controller before it sets the log level, and the race
	// detector flags that. The controller is started after hub.Parse
	// returns. See startExternalController.
	cfg, err := executor.ParseWithBytes(data)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	route.SetEmbedMode(true)
	// Drop the previous apply's proxies before hub.Parse. Their finalizers
	// call log.Debugln, and log.level has no lock, so that read races with
	// log.SetLevel inside ApplyConfig.
	releaseOutboundProxies()
	if err := hub.Parse(suppressExternalController(data)); err != nil {
		return fmt.Errorf("apply config: %w", err)
	}
	startExternalController(cfg.Controller.ExternalController, cfg.Controller.Secret, cfg.General.LogLevel == log.DEBUG)
	return nil
}

// suppressExternalController drops top-level external-controller lines so
// hub.Parse does not start the HTTP API. The file on disk is unchanged.
func suppressExternalController(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "external-controller") {
			continue
		}
		kept = append(kept, line)
	}
	return []byte(strings.Join(kept, "\n"))
}

func writeDefaultConfig(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, os.ErrExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}
	defer f.Close()
	if _, err := io.WriteString(f, DefaultConfig); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func releaseListeners() {
	listener.ReCreateHTTP(0, tunnel.Tunnel)
	listener.ReCreateSocks(0, tunnel.Tunnel)
	listener.ReCreateRedir(0, tunnel.Tunnel)
	listener.ReCreateTProxy(0, tunnel.Tunnel)
	listener.ReCreateMixed(0, tunnel.Tunnel)
	listener.Cleanup()
	mihomodns.ReCreateServer("", nil, nil)
	route.ReCreateServer(&route.Config{})
	executor.Shutdown()
}

// releaseOutboundProxies closes proxies on this goroutine and clears the
// auto-close finalizer. A later GC then has nothing left that logs.
func releaseOutboundProxies() {
	seen := map[C.Proxy]struct{}{}
	closeProxy := func(p C.Proxy) {
		if p == nil {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		_ = p.Close()
	}
	for _, p := range tunnel.Proxies() {
		closeProxy(p)
	}
	for _, pv := range tunnel.Providers() {
		if pv == nil {
			continue
		}
		for _, p := range pv.Proxies() {
			closeProxy(p)
		}
	}
	closeProviders()
}

func closeProviders() {
	for _, p := range tunnel.Providers() {
		closeIfPossible(p)
	}
	for _, p := range tunnel.RuleProviders() {
		closeIfPossible(p)
	}
}

func closeIfPossible(v any) {
	if c, ok := v.(io.Closer); ok && c != nil {
		_ = c.Close()
	}
}
