package ui

import (
	"strings"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

// Hooks are the app services a page calls. Nil hooks update the model only,
// which is what the headless tests do.
type Hooks struct {
	SetSystemProxy func(on bool) error
	SetTUN         func(on bool) error
	SetMode        func(mode string) error
	ImportURL      func(url string) error
	SavePrefs      func(Prefs) error
	OpenDir        func(kind string)
	RestartCore    func() error
	BrowseScript   func() (string, error)
	Copy           func(text string)
	Quit           func()
	TestSite       func(url string) error
}

// Model is the main window: a sidebar and the page it shows.
type Model struct {
	Router   *mygoui.Router
	Home     HomeState
	Cards    HomeCards
	Draft    HomeCards
	HomeEdit bool

	NetTab      int
	ChartSmooth bool
	Sites       []SiteTest
	IP          IPView
	IPShown     bool
	IPAuto      bool
	IPLeft      int

	Groups         []ProxyGroup
	Providers      []ProxyProvider
	Chain          ChainState
	ProxiesLoading bool
	Profiles       []ProfileItem
	CurrentProfile string
	ProfileURL     string
	ProfileBusy    bool
	Batch          bool
	ShowNextUpdate map[string]bool
	Form           ProfileForm
	DeleteUID      string
	Prefs          Prefs
	Notice         string
	Busy           string
	DNSAsk         bool
	ProfileHasDNS  bool
	ForceTray      bool
	Dialogs        dialogSet
	Hooks          Hooks
}

type dialogSet struct {
	Provider bool
	Theme    bool
	Layout   bool
	Misc     bool
	Hotkey   bool
	Backup   bool
	Ports    bool
	Sysproxy bool
	TUN      bool
	WebDAV   bool
	DNS      bool
	Runtime  bool
	Lite     bool
}

// NewModel returns a window opened on the home page.
func NewModel(home HomeState) *Model {
	if home.Mode == "" {
		home.Mode = "rule"
	}
	if home.VergeVersion == "" {
		home.VergeVersion = AppVersion
	}
	m := &Model{
		Router:         mygoui.NewRouter("/home"),
		Home:           home,
		Cards:          DefaultCards(),
		ChartSmooth:    true,
		Sites:          DefaultSites(),
		Prefs:          DefaultPrefs(),
		IPLeft:         30,
		ShowNextUpdate: map[string]bool{},
	}
	m.normalize()
	return m
}

func (m *Model) normalize() {
	if m == nil {
		return
	}
	switch m.Home.Mode {
	case "rule", "global", "direct":
	default:
		m.Home.Mode = "rule"
	}
	if m.Prefs.Language == "" {
		m.Prefs = DefaultPrefs()
	}
	if len(m.Sites) == 0 {
		m.Sites = DefaultSites()
	}
	if m.ShowNextUpdate == nil {
		m.ShowNextUpdate = map[string]bool{}
	}
}

// View builds one frame of the main window.
func (m *Model) View(c *mygoui.Context) {
	if m == nil {
		return
	}
	m.normalize()
	m.flushFilters(c.Now())
	dark := c.Theme().Dark
	switch m.Prefs.ThemeMode {
	case "light":
		dark = false
	case "dark":
		dark = true
	}
	pal := PaletteFor(dark)
	c.SetTheme(pal.Theme())
	c.Root().Background(pal.Paper)
	mygoui.Row(c).Fill().AlignItems(mygoui.Stretch).Children(func() {
		m.sidebar(c, pal)
		mygoui.Column(c).Grow(1).FillHeight().Background(pal.Content).Children(func() {
			m.Router.View(c, func(r *mygoui.Route) {
				path := m.Router.Path()
				switch path {
				case "/home":
					r.Title("首页")
					m.homePage(c, pal)
				case "/proxies":
					r.Title("代理组")
					m.proxyPage(c, pal)
				case "/profile":
					r.Title("订阅")
					m.profilePage(c, pal)
				case "/settings":
					r.Title("设置")
					m.settingsPage(c, pal)
				default:
					page, ok := pageByPath(path)
					if !ok {
						r.Title("未找到")
						PlaceholderPage(c, "未找到", "这个页面还没有。")
						return
					}
					r.Title(page.Title)
					PlaceholderPage(c, page.Title, page.Body)
				}
			})
		})
	})
	m.overlay(c, pal)
}

func (m *Model) flushFilters(now time.Time) {
	for i := range m.Groups {
		g := &m.Groups[i]
		if g.FilterOpen && Debounced(now, g.FilterChanged, FilterDebounce) {
			g.FilterApplied = strings.TrimSpace(g.FilterText)
			g.FilterChanged = time.Time{}
		}
	}
}

func (m *Model) sidebar(c *mygoui.Context, pal Palette) {
	page := m.Router.Path()
	mygoui.Column(c).Width(200).FillHeight().Background(pal.Paper).
		BorderWidth(0, 1, 0, 0).BorderColor(pal.Divider).Children(func() {
		mygoui.Column(c).Height(56).Padding(12, 16).Justify(mygoui.Center).Children(func() {
			mygoui.Text(c, "Clash Go").Bold().FontSize(16)
		})
		mygoui.Column(c).Grow(1).Padding(8, 12).Gap(4).Children(func() {
			for _, item := range Pages() {
				on := page == item.Path
				b := mygoui.ButtonBase(c.Key("nav-"+item.Path)).FillWidth().Height(36).Padding(0, 12).Radius(8).Label("nav-" + item.Nav)
				if on {
					b.Background(pal.SelectedFill())
					b.TextColor(pal.SidebarOn)
				} else {
					b.Background(mygoui.Transparent).TextColor(pal.Text)
				}
				b.AlignItems(mygoui.Center)
				nav := item
				b.Children(func() { mygoui.Text(c, nav.Nav).FontSize(14) })
				if b.Clicked() {
					m.Router.Push(nav.Path)
				}
			}
		})
		mygoui.Row(c).Height(60).Padding(0, 16).Gap(12).AlignItems(mygoui.Center).
			BorderWidth(1, 0, 0, 0).BorderColor(pal.Divider).Children(func() {
			mygoui.Column(c).Gap(2).Children(func() {
				mygoui.Text(c, "↑ "+FormatSpeed(m.Home.UpSpeed)).FontSize(12).TextColor(pal.TextMuted)
				mygoui.Text(c, "↓ "+FormatSpeed(m.Home.DownSpeed)).FontSize(12).TextColor(pal.TextMuted)
			})
			mygoui.Spacer(c).Grow(1)
			mygoui.Text(c, sidebarStatus(m.Home)).FontSize(11).TextColor(pal.TextMuted)
		})
	})
}

func sidebarStatus(s HomeState) string {
	if s.Running {
		return "运行中"
	}
	return "已停止"
}

func (m *Model) goPath(path string) {
	if path == "/" {
		path = "/home"
	}
	m.Router.Push(path)
}

func (m *Model) setMode(mode string) {
	if m.Home.Mode == mode {
		return
	}
	prev := m.Home.Mode
	m.Home.Mode = mode
	if m.Hooks.SetMode != nil {
		if err := m.Hooks.SetMode(mode); err != nil {
			m.Home.Mode = prev
			m.Notice = err.Error()
		}
	}
}

func (m *Model) setProxy(on bool) {
	if on && !m.Home.Running {
		m.Home.SystemProxy = false
		m.Prefs.SystemProxy = false
		m.Notice = "内核尚未就绪，无法开启系统代理"
		return
	}
	prev := m.Home.SystemProxy
	m.Home.SystemProxy = on
	m.Prefs.SystemProxy = on
	if m.Hooks.SetSystemProxy != nil {
		if err := m.Hooks.SetSystemProxy(on); err != nil {
			m.Home.SystemProxy = prev
			m.Prefs.SystemProxy = prev
			m.Notice = err.Error()
			return
		}
	}
}

func (m *Model) setTUN(on bool) {
	if on && !m.Home.TUNAvailable {
		m.Notice = "TUN 模式需要服务模式，请先安装服务"
		m.Home.TUN = false
		return
	}
	prev := m.Home.TUN
	m.Home.TUN = on
	m.Prefs.TUN = on
	if m.Hooks.SetTUN != nil {
		if err := m.Hooks.SetTUN(on); err != nil {
			m.Home.TUN = prev
			m.Prefs.TUN = prev
			m.Notice = err.Error()
		}
	}
}

func (m *Model) savePrefs() {
	if m.Hooks.SavePrefs != nil {
		if err := m.Hooks.SavePrefs(m.Prefs); err != nil {
			m.Notice = err.Error()
		}
	}
}

func (m *Model) copy(text string) {
	if m.Hooks.Copy != nil {
		m.Hooks.Copy(text)
	}
	m.Notice = "复制成功"
}

// CurrentItem is the active subscription, or nil.
func (m *Model) CurrentItem() *ProfileItem { return m.currentProfile() }

func (m *Model) currentProfile() *ProfileItem {
	for i := range m.Profiles {
		if m.Profiles[i].UID == m.CurrentProfile {
			return &m.Profiles[i]
		}
	}
	return nil
}

func (m *Model) subscriptionCount() int { return len(m.Profiles) }

func (m *Model) runMode() RunMode {
	if !m.Home.Running {
		return RunNotRunning
	}
	if m.Home.RunMode == "" {
		return RunSidecar
	}
	return m.Home.RunMode
}

func (m *Model) proxyAddr() string {
	if m.Home.ProxyAddr != "" {
		return m.Home.ProxyAddr
	}
	port := m.Home.MixedPort
	if port <= 0 {
		port = 7897
	}
	return "127.0.0.1:" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
