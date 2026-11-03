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
	// ActivateProfile makes uid the current subscription and reloads the
	// kernel so that profile's DNS overwrite takes effect.
	ActivateProfile func(uid string) error
	OpenDir         func(kind string)
	RestartCore     func() error
	BrowseScript    func() (string, error)
	Copy            func(text string)
	Quit            func()
	TestSite        func(url string) error
	// CloseAllConnections closes every active connection. A non-nil
	// error is kept off the page notice (Verge errorNotice: false).
	CloseAllConnections func() error
	// CloseConnection closes one active connection. Errors stay quiet.
	CloseConnection func(id string) error
	// RefreshConnections pulls a core snapshot into Active. now is the
	// frame clock, used to derive per-connection speeds.
	RefreshConnections func(now time.Time)
	RefreshRules       func()
	UpdateRuleProvider func(name string) error
	// RefreshLogs pulls kernel log lines recorded since the previous
	// call. The model drops them while clash log collection is paused.
	RefreshLogs func()
	// ClearLogs drops the kernel log buffer, the same role as Verge's
	// refreshGetClashLog(true). The page also clears its own rows.
	ClearLogs func()
	// RefreshUnlock reloads the unlock catalog. Nil keeps the page's
	// own merge of DefaultUnlockItems and the cached results.
	RefreshUnlock func()
	// CheckUnlock runs one service, or every service when name is empty.
	// onItem is called as each item finishes and may be nil. A non-nil
	// error is the detection failure notice. Tests inject this; the page
	// does not open a network connection itself.
	CheckUnlock func(name string, onItem func(UnlockItem)) ([]UnlockItem, error)
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

	Groups        []ProxyGroup
	Providers     []ProxyProvider
	Rules         []Rule
	RuleProviders []RuleProvider
	RuleQuery     string
	RulesList     mygoui.ListState
	Unlock        []UnlockItem
	PageVisible   bool
	// rulesOnPage and rulesPageVisible remember the last frame so a
	// refresh runs on enter and when the window is shown again.
	rulesOnPage      bool
	rulesPageVisible bool
	unlockOnPage     bool
	unlockAll        bool
	unlockLoading    map[string]bool
	Chain            ChainState
	ProxiesLoading   bool
	Profiles         []ProfileItem
	CurrentProfile   string
	ProfileURL       string
	ProfileBusy      bool
	Batch            bool
	ShowNextUpdate   map[string]bool
	Form             ProfileForm
	DeleteUID        string
	Prefs            Prefs
	Active           []Connection
	Closed           []Connection
	ConnUp           int64
	ConnDown         int64
	ConnLayout       string
	ConnOrder        string
	ConnQuery        string
	ConnSearch       SearchOpts
	ClashLog         ClashLog
	Logs             []LogItem
	LogQuery         string
	LogSearch        SearchOpts
	Notice           string
	noticeGen        uint64
	noticeShown      string
	noticeSentGen    uint64
	Busy             string
	DNSAsk           bool
	ProfileHasDNS    bool
	ForceTray        bool
	Dialogs          dialogSet
	Hooks            Hooks

	connTab      int
	connSelected int
	connSort     mygoui.SortOrder
	connDetailID string
	connColsOpen bool
	connHidden   map[string]bool
	connList     mygoui.ListState
	connTable    mygoui.ListState
	connScroll   mygoui.ScrollState
	connSynced   time.Time

	logList   mygoui.ListState
	logSynced time.Time
	logNext   uint64
}

type dialogSet struct {
	Provider     bool
	RuleProvider bool
	Theme        bool
	Layout       bool
	Misc         bool
	Hotkey       bool
	Backup       bool
	Ports        bool
	Sysproxy     bool
	TUN          bool
	WebDAV       bool
	DNS          bool
	Runtime      bool
	Lite         bool
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
		connSelected:   -1,
		PageVisible:    true,
		ClashLog:       DefaultClashLog(),
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
	if m.ConnLayout != "list" && m.ConnLayout != "table" {
		m.ConnLayout = "list"
	}
	if m.ConnOrder != "uploadSpeed" && m.ConnOrder != "downloadSpeed" {
		m.ConnOrder = "default"
	}
	if m.connHidden == nil {
		m.connHidden = map[string]bool{}
	}
	switch m.ClashLog.LogFilter {
	case "all", "debug", "info", "warn", "err":
	default:
		m.ClashLog.LogFilter = "all"
	}
	if m.ClashLog.LogOrder != "desc" {
		m.ClashLog.LogOrder = "asc"
	}
	if m.unlockLoading == nil {
		m.unlockLoading = map[string]bool{}
	}
}

// View builds one frame of the main window.
func (m *Model) View(c *mygoui.Context) {
	if m == nil {
		return
	}
	m.normalize()
	m.presentNotice(c)
	m.flushFilters(c.Now())
	m.syncRuleRefresh()
	m.syncUnlock()
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
				case "/connections":
					r.Title("连接")
					m.connectionsPage(c, pal)
				case "/rules":
					r.Title("规则")
					m.rulesPage(c, pal)
				case "/logs":
					r.Title("日志")
					m.logsPage(c, pal)
				case "/unlock":
					r.Title("解锁测试")
					m.unlockPage(c, pal)
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

func (m *Model) sidebarWidth() float32 {
	if m != nil && m.Prefs.CollapseNav {
		return sidebarCollapsedWidth
	}
	return sidebarExpandedWidth
}

func (m *Model) sidebar(c *mygoui.Context, pal Palette) {
	page := m.Router.Path()
	collapsed := m.Prefs.CollapseNav
	mygoui.Column(c).Width(m.sidebarWidth()).FillHeight().Background(pal.Paper).Label("sidebar").
		BorderWidth(0, 1, 0, 0).BorderColor(pal.Divider).Children(func() {
		m.logoRow(c, pal, collapsed)
		mygoui.Column(c).Grow(1).Padding(8, 10).Gap(4).Children(func() {
			for _, item := range Pages() {
				m.navItem(c, pal, item, page == item.Path, collapsed)
			}
		})
		if collapsed {
			return
		}
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

func (m *Model) navItem(c *mygoui.Context, pal Palette, item Page, on, collapsed bool) {
	fg := pal.Text
	if on {
		fg = pal.SidebarOn
	}
	if collapsed {
		mygoui.Row(c).FillWidth().Justify(mygoui.Center).Children(func() {
			b := mygoui.ButtonBase(c.Key("nav-"+item.Path)).Size(52, 52).Radius(12).Label("nav-" + item.Nav)
			if on {
				b.Background(pal.SelectedFill())
			} else {
				b.Background(mygoui.Transparent)
			}
			b.TextColor(fg).AlignItems(mygoui.Center).Justify(mygoui.Center)
			b.Children(func() {
				mygoui.Icon(c, navIcon(item.ID)).Size(22, 22)
			})
			if b.Clicked() {
				m.Router.Push(item.Path)
			}
		})
		return
	}
	b := mygoui.ButtonBase(c.Key("nav-"+item.Path)).FillWidth().Height(40).Padding(0, 12).Gap(10).Radius(8).Label("nav-" + item.Nav)
	if on {
		b.Background(pal.SelectedFill())
	} else {
		b.Background(mygoui.Transparent)
	}
	b.TextColor(fg).AlignItems(mygoui.Center)
	b.Children(func() {
		mygoui.Icon(c, navIcon(item.ID)).Size(20, 20).Shrink(0)
		mygoui.Text(c, item.Nav).FontSize(14)
	})
	if b.Clicked() {
		m.Router.Push(item.Path)
	}
}

func (m *Model) logoRow(c *mygoui.Context, pal Palette, collapsed bool) {
	if collapsed {
		// Collapsed mark uses margin 0 and sits centered in the 56×56 logo area.
		mygoui.Row(c).Height(logoCollapsedBox).AlignItems(mygoui.Center).Justify(mygoui.Center).Children(func() {
			mygoui.Column(c).Size(logoCollapsedBox, logoCollapsedBox).Label("logo-slot").AlignItems(mygoui.Center).Justify(mygoui.Center).Children(func() {
				mygoui.Image(c, brandMark).Size(brandMarkSize, brandMarkSize).Label("mygo")
			})
		})
		return
	}
	// Outer row is the sidebar logo slot. The inner row is Verge's 27px
	// line: the 36px mark uses margin -3/-3/5 so it sits against 「Clash Go」.
	mygoui.Row(c).Height(logoRowHeight).Padding(0, logoRowPaddingX).AlignItems(mygoui.Center).Children(func() {
		mygoui.Row(c).Height(logoInnerHeight).AlignItems(mygoui.Center).Children(func() {
			mygoui.Image(c, brandMark).Size(brandMarkSize, brandMarkSize).Shrink(0).
				Margin(brandMarkNudge, brandMarkMargin, 0, brandMarkNudge).Label("mygo")
			mygoui.Text(c, "Clash Go").Bold().FontSize(16).TextColor(pal.Text)
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
			m.SetNotice(err.Error())
		}
	}
}

func (m *Model) setProxy(on bool) {
	if on && !m.Home.Running {
		m.Home.SystemProxy = false
		m.Prefs.SystemProxy = false
		m.SetNotice("内核尚未就绪，无法开启系统代理")
		return
	}
	prev := m.Home.SystemProxy
	m.Home.SystemProxy = on
	m.Prefs.SystemProxy = on
	if m.Hooks.SetSystemProxy != nil {
		if err := m.Hooks.SetSystemProxy(on); err != nil {
			m.Home.SystemProxy = prev
			m.Prefs.SystemProxy = prev
			m.SetNotice(err.Error())
			return
		}
	}
}

func (m *Model) setTUN(on bool) {
	if on && !m.Home.TUNAvailable {
		m.SetNotice("TUN 模式需要服务模式，请先安装服务")
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
			m.SetNotice(err.Error())
		}
	}
}

func (m *Model) savePrefs() {
	if m.Hooks.SavePrefs != nil {
		if err := m.Hooks.SavePrefs(m.Prefs); err != nil {
			m.SetNotice(err.Error())
		}
	}
}

func (m *Model) copy(text string) {
	if m.Hooks.Copy != nil {
		m.Hooks.Copy(text)
	}
	m.SetNotice("复制成功")
}

// SetNotice records a page message and shows it as a transient toast.
// The same text can be sent again; an empty text clears the toast.
func (m *Model) SetNotice(text string) {
	if m == nil {
		return
	}
	m.Notice = text
	m.noticeGen++
}

func (m *Model) presentNotice(c *mygoui.Context) {
	if m.Notice == "" {
		if m.noticeShown != "" || m.noticeGen != m.noticeSentGen {
			c.CloseToast("clash-notice")
			m.noticeShown = ""
			m.noticeSentGen = m.noticeGen
		}
		return
	}
	if m.Notice == m.noticeShown && m.noticeGen == m.noticeSentGen {
		return
	}
	m.noticeShown = m.Notice
	m.noticeSentGen = m.noticeGen
	c.AddToast(mygoui.Toast{
		ID:      "clash-notice",
		Title:   m.Notice,
		Timeout: noticeTimeout(m.Notice),
	})
}

func noticeTimeout(text string) time.Duration {
	switch text {
	case "复制成功", "Verge 版本已复制到剪贴板":
		return time.Second
	default:
		return 4 * time.Second
	}
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
