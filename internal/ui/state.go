package ui

import "time"

// HomeCards is the home-settings checklist. Info has no card.
type HomeCards struct {
	Profile    bool
	Proxy      bool
	Network    bool
	Mode       bool
	Traffic    bool
	Test       bool
	IP         bool
	ClashInfo  bool
	SystemInfo bool
}

// DefaultCards matches Verge: every real card starts visible.
func DefaultCards() HomeCards {
	return HomeCards{
		Profile: true, Proxy: true, Network: true, Mode: true, Traffic: true,
		Test: true, IP: true, ClashInfo: true, SystemInfo: true,
	}
}

// HomeProfile is the subscription summary on the home card.
type HomeProfile struct {
	Name     string
	URL      string
	Home     string
	Updated  time.Time
	Used     int64
	Total    int64
	Expire   time.Time
	Updating bool
}

// CurrentNode is the node the home card summarizes.
type CurrentNode struct {
	Group      string
	Name       string
	Type       string
	Delay      int
	UDP        bool
	TFO        bool
	XUDP       bool
	MPTCP      bool
	SMUX       bool
	Groups     []string
	GroupNodes []ProxyNode
}

// ProxyNode is one member of a proxy group.
type ProxyNode struct {
	Name   string
	Type   string
	Delay  int
	Now    string
	UDP    bool
	TFO    bool
	XUDP   bool
	MPTCP  bool
	SMUX   bool
	Preset bool
}

// ProxyGroup is one collapsible group on the proxies page.
type ProxyGroup struct {
	Name          string
	Type          string
	Now           string
	Hidden        bool
	Icon          string
	Nodes         []ProxyNode
	Open          bool
	Sort          int
	ShowDetail    bool
	FilterOpen    bool
	FilterText    string
	FilterApplied string
	FilterChanged time.Time
	URLOpen       bool
	TestURL       string
}

// ProxyProvider is one entry in the provider dialog.
type ProxyProvider struct {
	Name     string
	Type     string
	Used     int64
	Total    int64
	Expire   string
	Updating bool
}

// Rule is one kernel rule. The line number is not stored: the page
// shows index+1 from the slice position and never writes it back.
type Rule struct {
	Type    string
	Payload string
	Proxy   string
}

// RuleProvider is one rule-set the rules page can update.
type RuleProvider struct {
	Name        string
	VehicleType string
	Behavior    string
	RuleCount   int
	UpdatedAt   time.Time
	Updating    bool
}

// ChainState is the chain-proxy panel.
type ChainState struct {
	On         bool
	Ask        bool
	Connected  bool
	Connecting bool
	Nodes      []ProxyNode
}

// SiteTest is one home website test.
type SiteTest struct {
	Name   string
	URL    string
	Status string // "", testing, ok, fail
}

// DefaultSites is the first-run Apple/GitHub/Google/YouTube list.
func DefaultSites() []SiteTest {
	return []SiteTest{
		{Name: "Apple", URL: "https://www.apple.com"},
		{Name: "GitHub", URL: "https://www.github.com"},
		{Name: "Google", URL: "https://www.google.com"},
		{Name: "YouTube", URL: "https://www.youtube.com"},
	}
}

// IPView is the IP card.
type IPView struct {
	Loaded   bool
	Failed   string
	IP       string
	Country  string
	Code     string
	ASN      string
	ISP      string
	Org      string
	Location string
	Timezone string
	Lon      string
	Lat      string
}

// ProfileItem is one subscription card.
type ProfileItem struct {
	UID         string
	Name        string
	Description string
	URL         string
	Type        string // remote, local
	Home        string
	Updated     time.Time
	Used        int64
	Total       int64
	Expire      time.Time
	Checked     bool
	Option      ProfileOption
}

// ProfileOption is the extra menu and the remote form.
type ProfileOption struct {
	UpdateInterval  int
	UseSystemProxy  bool
	UseClashProxy   bool
	AcceptInvalid   bool
	AllowAutoUpdate bool
	Rules           string
	Proxies         string
	Groups          string
	Merge           string
	Script          string
	UserAgent       string
	TimeoutSeconds  int
}

// ProfileForm is the create/edit dialog draft.
type ProfileForm struct {
	Open     bool
	Editing  bool
	UID      string
	Type     string
	Name     string
	Desc     string
	URL      string
	UA       string
	Timeout  int
	Interval int
	SysProxy bool
	Clash    bool
	Insecure bool
	Auto     bool
	Error    string
}

// Prefs is the settings page. Zero values are filled by DefaultPrefs.
type Prefs struct {
	SystemProxy    bool
	TUN            bool
	Language       string
	ThemeMode      string
	TrayEvent      string
	CopyEnv        string
	StartPage      string
	StartupScript  string
	AllowLan       bool
	IPv6           bool
	UnifiedDelay   bool
	LogLevel       string
	DNSOverwrite   bool
	ProxyColumns   int
	HoverDelay     int
	TrafficGraph   bool
	MemoryUsage    bool
	HoverNav       bool
	CollapseNav    bool
	GroupIcon      bool
	PauseOnBlur    bool
	AutoBackup     bool
	BackupEvery    int
	BackupOnChange bool
	SilentStart    bool
	AutoLaunch     bool
	Lightweight    bool
	InlineModes    bool
	ProxyMenu      string
	WebDAVURL      string
	WebDAVUser     string
	WebDAVPassword string
	ThemePrimary   string
	ThemeSecondary string
	ThemeText      string
	ThemeMuted     string
	ThemeInfo      string
	ThemeWarning   string
	ThemeError     string
	ThemeSuccess   string
	FontFamily     string
	MixedPort      int
	SocksPort      int
	SocksOn        bool
	HTTPPort       int
	HTTPOn         bool
	RedirPort      int
	RedirOn        bool
	TProxyPort     int
	TProxyOn       bool
	HotkeyDash     string
	HotkeyRule     string
	HotkeyGlobal   string
	HotkeyDirect   string
	HotkeyProxy    string
	HotkeyTUN      string
	HotkeyLite     string
	HotkeyProfiles string
}

// DefaultPrefs matches the settings spec defaults the page must show.
func DefaultPrefs() Prefs {
	return Prefs{
		Language:       "zh",
		ThemeMode:      "system",
		TrayEvent:      "main_window",
		CopyEnv:        "bash",
		StartPage:      "/",
		LogLevel:       "info",
		IPv6:           true,
		UnifiedDelay:   true,
		ProxyColumns:   AutoProxyColumns,
		HoverDelay:     HoverJumpDelay,
		TrafficGraph:   true,
		MemoryUsage:    true,
		HoverNav:       true,
		GroupIcon:      true,
		PauseOnBlur:    true,
		BackupEvery:    24,
		BackupOnChange: true,
		ProxyMenu:      "default",
		MixedPort:      7897,
		SocksPort:      7898,
		HTTPPort:       7899,
		RedirPort:      7895,
		TProxyPort:     7896,
	}
}

// RunMode is the system-info badge.
type RunMode string

const (
	RunNotRunning RunMode = "not-running"
	RunSidecar    RunMode = "sidecar"
	RunService    RunMode = "service"
	RunAdmin      RunMode = "admin"
	RunAdminSvc   RunMode = "admin-service"
)

// RunModeLabel is the Chinese badge.
func RunModeLabel(m RunMode) string {
	switch m {
	case RunNotRunning:
		return "内核未运行"
	case RunAdminSvc:
		return "管理员 + 服务模式"
	case RunAdmin:
		return "管理员模式"
	case RunSidecar:
		return "用户模式"
	default:
		return "服务模式"
	}
}

// LanguageChoice is one endonym in the language menu.
type LanguageChoice struct {
	ID   string
	Name string
}

// Languages are the settings dropdown endonyms, not translated sentences.
func Languages() []LanguageChoice {
	return []LanguageChoice{
		{"en", "English"},
		{"ru", "Русский"},
		{"zh", "简体中文"},
		{"fa", "فارسی"},
		{"tt", "Татар"},
		{"id", "Bahasa Indonesia"},
		{"ar", "العربية"},
		{"ko", "한국어"},
		{"tr", "Türkçe"},
		{"de", "Deutsch"},
		{"es", "Español"},
		{"ja", "日本語"},
		{"zh-TW", "繁體中文"},
		{"pt", "Português"},
	}
}

// LanguageName returns the endonym for a code, or the code itself.
func LanguageName(id string) string {
	for _, l := range Languages() {
		if l.ID == id {
			return l.Name
		}
	}
	if id == "" {
		return "简体中文"
	}
	return id
}

// LanguageID maps an endonym back to a code.
func LanguageID(name string) string {
	for _, l := range Languages() {
		if l.Name == name {
			return l.ID
		}
	}
	return name
}

// TrayEvents are the macOS tray-click choices. Linux hides the row.
func TrayEvents() []LanguageChoice {
	return []LanguageChoice{
		{"main_window", "显示主窗口"},
		{"tray_menu", "显示托盘菜单"},
		{"system_proxy", "系统代理"},
		{"tun", "虚拟网卡模式"},
		{"disable", "禁用"},
	}
}

// CopyEnvTypes are the environment-variable shells.
func CopyEnvTypes() []string { return []string{"bash", "fish", "nushell", "cmd", "powershell"} }

// StartPages are the eight navigation targets, using the sidebar labels.
func StartPages() []LanguageChoice {
	var out []LanguageChoice
	for _, p := range Pages() {
		out = append(out, LanguageChoice{ID: p.Path, Name: p.Nav})
	}
	return out
}
