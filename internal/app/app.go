// Package app wires the kernel, settings and the native window.
package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/profile"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
	"github.com/charlzyx/clash-go/internal/ui"

	"github.com/egoist/mygo"
)

// Options selects the pieces a test can replace. Zero values get the
// user config directory, a file store, the macOS system proxy on darwin
// (an in-memory proxy elsewhere), and a TUN controller that only
// remembers the switch.
type Options struct {
	Dirs  appdir.Provider
	Store appconfig.Store
	Proxy sysproxy.Controller
	TUN   tun.Controller
	Core  *core.Core
}

// App is one running Clash Go process.
type App struct {
	dirs  appdir.Provider
	store appconfig.Store
	proxy sysproxy.Controller
	tun   tun.Controller
	core  *core.Core

	settings  appconfig.Settings
	model     *ui.Model
	dir       string
	profiles  *profile.Store
	started   time.Time
	logSeq    uint64
	applied   clashStamp
	appliedOK bool

	win  *mygo.Window
	tray *mygo.Tray
}

// New checks options and fills the M1 defaults. It does not start the
// kernel or the window.
func New(opt Options) *App {
	if opt.Dirs == nil {
		opt.Dirs = appdir.Default{}
	}
	if opt.Proxy == nil {
		opt.Proxy = sysproxy.Platform()
	}
	if opt.TUN == nil {
		opt.TUN = tun.NewMemory()
	}
	if opt.Core == nil {
		opt.Core = core.New()
	}
	return &App{
		dirs:  opt.Dirs,
		store: opt.Store,
		proxy: opt.Proxy,
		tun:   opt.TUN,
		core:  opt.Core,
	}
}

// Prepare loads settings, starts the kernel in dir, and builds the home
// snapshot. The window is created separately so tests can stop here.
func (a *App) Prepare() error {
	dir, err := a.dirs.ConfigDir()
	if err != nil {
		return err
	}
	a.dir = dir
	if a.store == nil {
		a.store = appconfig.NewFileStore(dir)
	}
	settings, err := a.store.Load()
	if err != nil {
		return err
	}
	a.settings = settings
	a.profiles = profile.NewStore(dir)
	a.model = ui.NewModel(ui.HomeState{ConfigDir: dir, Version: core.Version()})
	a.applySettings()
	a.loadProfiles()
	if err := a.installCurrentProfile(); err != nil {
		a.model.SetNotice(err.Error())
	}
	if err := a.core.Start(dir); err != nil {
		a.model.Home.Version = core.Version()
		a.model.Home.Running = false
		a.model.Home.ConfigDir = dir
		a.model.Home.Error = err.Error()
		a.model.Home.SystemProxy = settings.SystemProxy
		a.model.Home.TUN = settings.TUN
		a.bindHooks()
		if settings.SystemProxy {
			a.settings.SystemProxy = false
			a.model.Home.SystemProxy = false
			a.model.Prefs.SystemProxy = false
			if a.store != nil {
				_ = a.store.Save(a.settings)
			}
			a.model.SetNotice("内核尚未就绪，系统代理保持关闭")
		}
		return nil
	}
	a.started = time.Now()
	a.loadGroups()
	a.loadRules()
	a.bindHooks()
	a.refreshHome()
	a.markApplied()
	a.reconcileSystemProxy()
	return nil
}

func (a *App) applySettings() {
	prefs := ui.DefaultPrefs()
	if a.settings.Prefs != "" {
		_ = json.Unmarshal([]byte(a.settings.Prefs), &prefs)
	}
	prefs.SystemProxy = a.settings.SystemProxy
	prefs.TUN = a.settings.TUN
	if a.settings.TrayEvent != "" {
		prefs.TrayEvent = a.settings.TrayEvent
	}
	if a.settings.ProxyMenu != "" {
		prefs.ProxyMenu = a.settings.ProxyMenu
	}
	prefs.InlineModes = a.settings.InlineModes
	prefs.Lightweight = a.settings.Lightweight
	a.model.Prefs = prefs
	a.model.Home.SystemProxy = prefs.SystemProxy
	a.model.Home.TUN = prefs.TUN
	a.model.Home.AutoLaunch = prefs.AutoLaunch
}

func (a *App) loadProfiles() {
	if a.profiles == nil {
		return
	}
	file, err := a.profiles.Load()
	if err != nil {
		a.model.SetNotice(err.Error())
		return
	}
	a.model.Profiles = nil
	a.model.CurrentProfile = file.Current
	for _, it := range file.Items {
		a.model.Profiles = append(a.model.Profiles, ui.ProfileItem{
			UID: it.UID, Name: it.Name, Description: it.Description, URL: it.URL, Type: it.Type,
			Home: it.Home, Updated: it.Updated, Used: it.Used, Total: it.Total, Expire: it.Expire,
			Option: ui.ProfileOption{
				UpdateInterval: it.Option.UpdateInterval, UseSystemProxy: it.Option.UseSystemProxy,
				UseClashProxy: it.Option.UseClashProxy, AcceptInvalid: it.Option.AcceptInvalid,
				AllowAutoUpdate: it.Option.AllowAutoUpdate, Rules: it.Option.Rules, Proxies: it.Option.Proxies,
				Groups: it.Option.Groups, Merge: it.Option.Merge, Script: it.Option.Script,
				UserAgent: it.Option.UserAgent, TimeoutSeconds: it.Option.TimeoutSeconds,
				DNSOverwrite: it.Option.DNSOverwrite,
			},
		})
	}
	if cur := a.model.CurrentItem(); cur != nil {
		a.model.Home.Profile = &ui.HomeProfile{Name: cur.Name, URL: cur.URL, Home: cur.Home, Updated: cur.Updated, Used: cur.Used, Total: cur.Total, Expire: cur.Expire}
		a.model.Prefs.DNSOverwrite = cur.Option.DNSOverwrite
	}
	a.model.ProfileHasDNS = profile.HasTopLevelDNS(a.currentProfileBody())
}

// SetSystemProxy asks the proxy controller to enable or disable the
// system proxy and stores the result. It does not read the kernel
// singleton from the UI; the endpoint comes from the last status.
func (a *App) SetSystemProxy(on bool) error {
	endpoint := a.proxyEndpoint()
	next := sysproxy.Settings{Enabled: on, HTTP: endpoint, SOCKS: endpoint}
	if err := a.proxy.Set(next); err != nil {
		return err
	}
	prev := a.settings.SystemProxy
	a.settings.SystemProxy = on
	if err := a.store.Save(a.settings); err != nil {
		a.settings.SystemProxy = prev
		_ = a.proxy.Set(sysproxy.Settings{Enabled: prev, HTTP: endpoint, SOCKS: endpoint})
		return err
	}
	a.refreshHome()
	return nil
}

// SetTUN asks the TUN controller to enable or disable TUN and stores it.
func (a *App) SetTUN(on bool) error {
	if err := a.tun.Set(on); err != nil {
		return err
	}
	prev := a.settings.TUN
	a.settings.TUN = on
	if err := a.store.Save(a.settings); err != nil {
		a.settings.TUN = prev
		_ = a.tun.Set(prev)
		return err
	}
	a.refreshHome()
	return nil
}

// Home is the snapshot the window shows.
func (a *App) Home() ui.HomeState {
	if a == nil || a.model == nil {
		return ui.HomeState{}
	}
	return a.model.Home
}

// Model is the window model. Nil before Prepare.
func (a *App) Model() *ui.Model { return a.model }

func (a *App) refreshHome() {
	if a.model == nil {
		return
	}
	st := a.core.Status()
	h := a.model.Home
	h.Version = st.Version
	h.Running = st.Running
	if st.MixedPort > 0 {
		h.MixedPort = st.MixedPort
	}
	if st.Mode != "" {
		h.Mode = st.Mode
	}
	h.ConfigDir = a.dir
	h.SystemProxy = a.settings.SystemProxy
	h.TUN = a.settings.TUN
	if st.Running && !a.started.IsZero() {
		h.Uptime = time.Since(a.started)
	}
	a.model.Home = h
	a.model.Prefs.SystemProxy = a.settings.SystemProxy
	a.model.Prefs.TUN = a.settings.TUN
	if a.win != nil {
		a.win.Invalidate()
	}
	if a.tray != nil {
		a.tray.SetMenu(BuildTrayMenu(a.traySnapshot(), a.trayActions()))
	}
}

func (a *App) installCurrentProfile() error {
	return a.writeRuntime(false)
}

// ActivateProfile makes uid the current subscription and reloads the kernel
// with that profile's DNS overwrite plus the clash fields from prefs.
func (a *App) ActivateProfile(uid string) error {
	if a.profiles == nil || a.dir == "" {
		return nil
	}
	file, err := a.profiles.Load()
	if err != nil {
		return err
	}
	file.Current = uid
	if err := a.profiles.Save(file); err != nil {
		return err
	}
	a.loadProfiles()
	if err := a.writeRuntime(false); err != nil {
		return err
	}
	if a.core != nil && a.core.Status().Running {
		if err := a.core.Reload(); err != nil {
			return err
		}
		a.loadGroups()
		a.loadRules()
		a.refreshHome()
	}
	a.markApplied()
	return nil
}

func (a *App) loadGroups() {
	if a.model == nil || a.dir == "" {
		return
	}
	body, err := os.ReadFile(filepath.Join(a.dir, "config.yaml"))
	if err != nil {
		return
	}
	parsed, err := profile.ParseClashYAML(body)
	if err != nil {
		a.model.SetNotice(err.Error())
		return
	}
	if parsed.Rules > 0 {
		a.model.Home.Rules = parsed.Rules
	}
	groups := make([]ui.ProxyGroup, 0, len(parsed.Groups))
	for _, g := range parsed.Groups {
		gg := ui.ProxyGroup{Name: g.Name, Type: g.Type, Now: g.Now, Open: true, ShowDetail: true}
		for _, n := range g.Proxies {
			gg.Nodes = append(gg.Nodes, ui.ProxyNode{Name: n.Name, Type: n.Type, Delay: -1})
		}
		groups = append(groups, gg)
	}
	if len(groups) > 0 {
		a.model.Groups = groups
		a.model.Home.Node.Group = groups[0].Name
		a.model.Home.Node.Groups = nil
		for _, g := range groups {
			a.model.Home.Node.Groups = append(a.model.Home.Node.Groups, g.Name)
			if g.Now != "" && a.model.Home.Node.Name == "" {
				a.model.Home.Node.Name = g.Now
				a.model.Home.Node.Group = g.Name
			}
		}
	}
}

func (a *App) bindHooks() {
	if a.model == nil {
		return
	}
	a.model.Hooks = ui.Hooks{
		SetSystemProxy:      a.SetSystemProxy,
		SetTUN:              a.SetTUN,
		SetMode:             a.SetMode,
		ImportURL:           a.ImportURL,
		SavePrefs:           a.savePrefs,
		ActivateProfile:     a.ActivateProfile,
		OpenDir:             a.openDir,
		RestartCore:         a.restartCore,
		Copy:                writeClipboard,
		Quit:                func() { mygo.App.Quit() },
		CloseAllConnections: a.closeAllConnections,
		CloseConnection:     a.closeConnection,
		RefreshConnections:  a.refreshConnections,
		RefreshRules:        a.loadRules,
		UpdateRuleProvider:  a.updateRuleProvider,
		RefreshLogs:         a.refreshLogs,
		ClearLogs:           a.clearLogs,
	}
}

// loadRules copies the kernel rule list and rule providers into the window.
// A stopped core leaves whatever the page already shows.
func (a *App) loadRules() {
	if a.model == nil || a.core == nil || !a.core.Status().Running {
		return
	}
	raw := a.core.Rules()
	rules := make([]ui.Rule, 0, len(raw))
	for _, rule := range raw {
		rules = append(rules, ui.Rule{Type: rule.Type, Payload: rule.Payload, Proxy: rule.Proxy})
	}
	provs := a.core.RuleProviders()
	out := make([]ui.RuleProvider, 0, len(provs))
	for _, p := range provs {
		out = append(out, ui.RuleProvider{
			Name: p.Name, VehicleType: p.VehicleType, Behavior: p.Behavior,
			RuleCount: p.RuleCount, UpdatedAt: p.UpdatedAt,
		})
	}
	a.model.Rules = rules
	a.model.RuleProviders = out
	a.model.Home.Rules = len(rules)
	if a.win != nil {
		a.win.Invalidate()
	}
}

func (a *App) updateRuleProvider(name string) error {
	if a.core == nil {
		return nil
	}
	return a.core.UpdateRuleProvider(name)
}

func (a *App) setPageVisible(on bool) {
	if a.model == nil || a.model.PageVisible == on {
		return
	}
	a.model.PageVisible = on
	if a.win != nil {
		a.win.Invalidate()
	}
}

func (a *App) refreshLogs() {
	if a.model == nil || a.core == nil {
		return
	}
	lines, next, ok := a.core.LogsSince(a.logSeq)
	if !ok {
		return
	}
	a.logSeq = next
	if len(lines) == 0 {
		return
	}
	items := make([]ui.LogItem, len(lines))
	for i, line := range lines {
		items[i] = ui.LogItem{
			Time:    line.Time.Format("01-02 15:04:05"),
			Type:    line.Level,
			Payload: line.Payload,
		}
	}
	a.model.AppendLogs(items)
}

func (a *App) clearLogs() {
	if a.core != nil {
		a.core.ClearLogs()
	}
	if a.model != nil {
		a.model.Logs = nil
	}
}

func (a *App) refreshConnections(now time.Time) {
	if a.model == nil || a.core == nil {
		return
	}
	up, down, items, ok := a.core.ConnectionSnapshot()
	if !ok {
		return
	}
	view := make([]ui.Connection, len(items))
	for i, it := range items {
		view[i] = ui.Connection{
			ID: it.ID, Host: it.Host, DestinationIP: it.DestinationIP, DestPort: it.DestPort,
			RemoteDest: it.RemoteDestination, Process: it.Process, ProcessPath: it.ProcessPath,
			Network: it.Network, Type: it.Type, SourceIP: it.SourceIP, SourcePort: it.SourcePort,
			Rule: it.Rule, RulePayload: it.RulePayload, Chains: it.Chains,
			Upload: it.Upload, Download: it.Download, Start: it.Start,
		}
	}
	a.model.ApplyConnections(view, up, down, now)
}

func (a *App) closeAllConnections() error {
	if a.core != nil && a.core.Status().Running {
		if err := a.core.CloseAllConnections(); err != nil {
			return err
		}
	}
	a.refreshConnections(time.Now())
	return nil
}

func (a *App) closeConnection(id string) error {
	if a.core != nil && a.core.Status().Running {
		if err := a.core.CloseConnection(id); err != nil {
			return err
		}
	}
	a.refreshConnections(time.Now())
	return nil
}

// SetMode writes the profile mode and reloads the kernel when it is running.
func (a *App) SetMode(mode string) error {
	if a.dir == "" {
		return nil
	}
	path := filepath.Join(a.dir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, profile.SetMode(data, mode), 0o644); err != nil {
		return err
	}
	if a.core != nil && a.core.Status().Running {
		if err := a.core.Reload(); err != nil {
			return err
		}
	}
	if a.model != nil {
		a.model.Home.Mode = mode
	}
	return nil
}

// ImportURL stores a remote subscription and, when it is the current one,
// installs it as the kernel profile.
func (a *App) ImportURL(raw string) error {
	if a.profiles == nil {
		a.profiles = profile.NewStore(a.dir)
	}
	if _, err := a.profiles.ImportRemote(raw, fetchProfile); err != nil {
		return err
	}
	a.loadProfiles()
	if err := a.installCurrentProfile(); err != nil {
		return err
	}
	if a.core != nil && a.core.Status().Running {
		if err := a.core.Reload(); err != nil {
			return err
		}
	}
	a.loadGroups()
	a.loadRules()
	a.markApplied()
	return nil
}

func (a *App) savePrefs(p ui.Prefs) error {
	if a.model != nil {
		a.model.Prefs = p
	}
	if err := a.persistSettings(p); err != nil {
		return err
	}
	if err := a.persistProfileDNS(p.DNSOverwrite); err != nil {
		return err
	}
	next := a.clashStamp(p)
	if a.appliedOK && next == a.applied {
		return nil
	}
	if err := a.writeRuntime(true); err != nil {
		return err
	}
	if a.core != nil && a.core.Status().Running {
		if err := a.core.Reload(); err != nil {
			return err
		}
		a.loadGroups()
		a.loadRules()
		a.refreshHome()
	}
	a.applied = next
	a.appliedOK = true
	return nil
}

func (a *App) persistSettings(p ui.Prefs) error {
	a.settings.SystemProxy = p.SystemProxy
	a.settings.TUN = p.TUN
	a.settings.TrayEvent = p.TrayEvent
	a.settings.InlineModes = p.InlineModes
	a.settings.ProxyMenu = p.ProxyMenu
	a.settings.Lightweight = p.Lightweight
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	a.settings.Prefs = string(b)
	if a.store == nil {
		return nil
	}
	return a.store.Save(a.settings)
}

// clashStamp is the subset of prefs that changes config.yaml. Theme, language
// and the other Verge fields are not in it, so saving those does not reload.
type clashStamp struct {
	AllowLan     bool
	IPv6         bool
	UnifiedDelay bool
	DNS          bool
	LogLevel     string
	Mixed        int
	Socks        int
	SocksOn      bool
	HTTP         int
	HTTPOn       bool
	Redir        int
	RedirOn      bool
	TProxy       int
	TProxyOn     bool
	Profile      string
}

func (a *App) clashStamp(p ui.Prefs) clashStamp {
	dns := p.DNSOverwrite
	uid := ""
	if a.model != nil {
		uid = a.model.CurrentProfile
		if cur := a.model.CurrentItem(); cur != nil {
			dns = cur.Option.DNSOverwrite
		}
	}
	return clashStamp{
		AllowLan: p.AllowLan, IPv6: p.IPv6, UnifiedDelay: p.UnifiedDelay, DNS: dns,
		LogLevel: profile.KernelLogLevel(p.LogLevel),
		Mixed:    p.MixedPort, Socks: p.SocksPort, SocksOn: p.SocksOn,
		HTTP: p.HTTPPort, HTTPOn: p.HTTPOn,
		Redir: p.RedirPort, RedirOn: p.RedirOn,
		TProxy: p.TProxyPort, TProxyOn: p.TProxyOn,
		Profile: uid,
	}
}

func (a *App) markApplied() {
	if a.model == nil {
		return
	}
	a.applied = a.clashStamp(a.model.Prefs)
	a.appliedOK = true
}

func (a *App) clashPatch(p ui.Prefs) profile.ClashPatch {
	dns := p.DNSOverwrite
	if a.model != nil {
		if cur := a.model.CurrentItem(); cur != nil {
			dns = cur.Option.DNSOverwrite
		}
	}
	return profile.ClashPatch{
		AllowLan: p.AllowLan, IPv6: p.IPv6, UnifiedDelay: p.UnifiedDelay,
		LogLevel: p.LogLevel, MixedPort: p.MixedPort,
		SocksPort: p.SocksPort, SocksOn: p.SocksOn,
		HTTPPort: p.HTTPPort, HTTPOn: p.HTTPOn,
		RedirPort: p.RedirPort, RedirOn: p.RedirOn, ApplyRedir: runtime.GOOS != "windows",
		TProxyPort: p.TProxyPort, TProxyOn: p.TProxyOn, ApplyTProxy: runtime.GOOS == "linux",
		DNSOverwrite: dns,
	}
}

// writeRuntime installs the current subscription (or the existing config,
// or the built-in default) and then stamps clash prefs on top. preserveMode
// keeps the mode the user picked; a profile install lets the file's own
// mode through.
func (a *App) writeRuntime(preserveMode bool) error {
	if a.dir == "" {
		return nil
	}
	body := a.runtimeBase()
	mode := ""
	if preserveMode && a.model != nil {
		mode = a.model.Home.Mode
	}
	prefs := ui.DefaultPrefs()
	if a.model != nil {
		prefs = a.model.Prefs
	}
	body = profile.ApplyClash(body, a.clashPatch(prefs))
	if mode != "" {
		body = profile.SetMode(body, mode)
	}
	if err := os.MkdirAll(a.dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(a.dir, "config.yaml"), body, 0o644)
}

func (a *App) runtimeBase() []byte {
	if body := a.currentProfileBody(); len(body) > 0 {
		return body
	}
	if a.dir != "" {
		if data, err := os.ReadFile(filepath.Join(a.dir, "config.yaml")); err == nil && len(bytes.TrimSpace(data)) > 0 {
			return data
		}
	}
	return []byte(core.DefaultConfig)
}

func (a *App) currentProfileBody() []byte {
	if a.profiles == nil || a.model == nil || a.dir == "" {
		return nil
	}
	cur := a.model.CurrentItem()
	if cur == nil {
		return nil
	}
	file, err := a.profiles.Load()
	if err != nil {
		return nil
	}
	var rel string
	for _, it := range file.Items {
		if it.UID == cur.UID {
			rel = it.File
		}
	}
	if rel == "" {
		return nil
	}
	body, err := os.ReadFile(filepath.Join(a.dir, rel))
	if err != nil {
		return nil
	}
	return body
}

func (a *App) persistProfileDNS(on bool) error {
	if a.model != nil {
		a.model.Prefs.DNSOverwrite = on
		if cur := a.model.CurrentItem(); cur != nil {
			cur.Option.DNSOverwrite = on
		}
	}
	if a.profiles == nil || a.model == nil || a.model.CurrentProfile == "" {
		return nil
	}
	file, err := a.profiles.Load()
	if err != nil {
		return err
	}
	changed := false
	for i := range file.Items {
		if file.Items[i].UID != a.model.CurrentProfile {
			continue
		}
		if file.Items[i].Option.DNSOverwrite == on {
			return nil
		}
		file.Items[i].Option.DNSOverwrite = on
		changed = true
	}
	if !changed {
		return nil
	}
	return a.profiles.Save(file)
}

// reconcileSystemProxy applies a saved "on" to the real controller. A
// platform that cannot do it clears the flag instead of leaving the switch
// on.
func (a *App) reconcileSystemProxy() {
	if a.model == nil || !a.settings.SystemProxy {
		return
	}
	if err := a.SetSystemProxy(true); err != nil {
		a.settings.SystemProxy = false
		a.model.Home.SystemProxy = false
		a.model.Prefs.SystemProxy = false
		if a.store != nil {
			_ = a.store.Save(a.settings)
		}
		a.model.SetNotice("系统代理未启用：" + err.Error())
	}
}

func (a *App) openDir(kind string) {
	path := a.dir
	switch kind {
	case "logs", "app-log", "core-log":
		path = filepath.Join(a.dir, "logs")
		_ = os.MkdirAll(path, 0o755)
	}
	if path == "" {
		return
	}
	_ = openPath(path)
}

func (a *App) restartCore() error {
	if a.core == nil || a.dir == "" {
		return nil
	}
	_ = a.core.Stop()
	if err := a.core.Start(a.dir); err != nil {
		return err
	}
	a.started = time.Now()
	a.loadGroups()
	a.loadRules()
	a.refreshHome()
	return nil
}

func fetchProfile(raw string) ([]byte, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(raw)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, io.ErrUnexpectedEOF
	}
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}

func (a *App) proxyEndpoint() string {
	port := 7897
	if a.core != nil {
		if st := a.core.Status(); st.MixedPort > 0 {
			port = st.MixedPort
		}
	}
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
}
