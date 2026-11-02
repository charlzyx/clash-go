// Package app wires the kernel, settings and the native window.
package app

import (
	"net"
	"strconv"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
	"github.com/charlzyx/clash-go/internal/ui"

	"github.com/egoist/mygo"
)

// Options selects the pieces a test can replace. Zero values get the
// M1 defaults: the user config directory, a file store, and controllers
// that remember switches without touching the OS.
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

	settings appconfig.Settings
	model    *ui.Model
	dir      string

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
		opt.Proxy = sysproxy.NewMemory()
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
	a.model = ui.NewModel(ui.HomeState{})
	if err := a.core.Start(dir); err != nil {
		a.model.Home = ui.HomeState{
			Version:     core.Version(),
			Running:     false,
			ConfigDir:   dir,
			Error:       err.Error(),
			SystemProxy: settings.SystemProxy,
			TUN:         settings.TUN,
		}
		return nil
	}
	a.refreshHome()
	return nil
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
	a.model.Home = ui.HomeState{
		Version:     st.Version,
		Running:     st.Running,
		MixedPort:   st.MixedPort,
		Mode:        st.Mode,
		ConfigDir:   a.dir,
		SystemProxy: a.settings.SystemProxy,
		TUN:         a.settings.TUN,
	}
	if a.win != nil {
		a.win.Invalidate()
	}
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
