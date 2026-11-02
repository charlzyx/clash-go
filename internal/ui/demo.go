package ui

import (
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

// DemoModel is the populated window the M2 screenshots render.
func DemoModel() *Model {
	m := NewModel(HomeState{
		Version:      "1.19.33",
		Running:      true,
		MixedPort:    7897,
		Mode:         "rule",
		ConfigDir:    "~/Library/Application Support/app.clashgo.desktop",
		SystemProxy:  true,
		TUN:          false,
		TUNAvailable: false,
		Rules:        128,
		Uptime:       time.Hour + 2*time.Minute + 3*time.Second,
		UpSpeed:      18 * 1024,
		DownSpeed:    240 * 1024,
		Connections:  12,
		TotalUp:      800 * 1024 * 1024,
		TotalDown:    6 * 1024 * 1024 * 1024,
		Memory:       48 * 1024 * 1024,
		OS:           "macOS",
		OSVersion:    "15.1",
		RunMode:      RunSidecar,
		VergeVersion: AppVersion,
		Profile: &HomeProfile{
			Name:    "示例订阅",
			URL:     "https://example.com/sub",
			Home:    "https://example.com",
			Updated: time.Date(2026, 11, 2, 9, 30, 0, 0, time.UTC),
			Used:    12 * 1024 * 1024 * 1024,
			Total:   100 * 1024 * 1024 * 1024,
			Expire:  time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		Node: CurrentNode{
			Group:  "Proxy",
			Name:   "香港 01",
			Type:   "ss",
			Delay:  128,
			UDP:    true,
			TFO:    true,
			Groups: []string{"Proxy", "Auto"},
		},
	})
	m.ForceTray = true
	m.IP = IPView{IP: "203.0.113.10", Country: "中国", Code: "CN", ASN: "AS4134", ISP: "CHINANET", Org: "CHINANET", Location: "Shanghai, Shanghai", Timezone: "Asia/Shanghai", Lon: "121.47", Lat: "31.23"}
	m.Sites = DefaultSites()
	hk := []ProxyNode{
		{Name: "香港 01", Type: "ss", Delay: 128, UDP: true, TFO: true},
		{Name: "香港 02", Type: "vmess", Delay: 280},
		{Name: "台湾 01", Type: "trojan", Delay: 460},
		{Name: "美国 01", Type: "ss", Delay: -1},
	}
	m.Groups = []ProxyGroup{
		{Name: "Proxy", Type: "Selector", Now: "香港 01", Open: true, ShowDetail: true, Nodes: hk},
		{Name: "Auto", Type: "URLTest", Now: "香港 02", Open: true, ShowDetail: true, Nodes: hk},
	}
	m.Providers = []ProxyProvider{{Name: "airport", Type: "http", Used: 12 << 30, Total: 100 << 30, Expire: "2026-12-31"}}
	m.Profiles = []ProfileItem{
		{UID: "p1", Name: "示例订阅", URL: "https://example.com/sub", Type: "remote", Home: "https://example.com", Updated: time.Date(2026, 11, 2, 9, 30, 0, 0, time.UTC), Used: 12 << 30, Total: 100 << 30, Expire: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), Option: ProfileOption{UpdateInterval: 1440, AllowAutoUpdate: true, Rules: "x"}},
		{UID: "p2", Name: "本地配置", Type: "local", Description: "手动维护的规则", Updated: time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)},
	}
	m.CurrentProfile = "p1"
	m.Prefs.SystemProxy = true
	m.ConnDown = 6 * 1024 * 1024 * 1024
	m.ConnUp = 800 * 1024 * 1024
	m.Active = []Connection{
		{
			ID: "c1", Host: "example.com", DestinationIP: "203.0.113.10", DestPort: 443,
			Process: "Safari", ProcessPath: "/Applications/Safari.app", Network: "tcp", Type: "HTTPS",
			SourceIP: "127.0.0.1", SourcePort: 52314, Rule: "DOMAIN-SUFFIX", RulePayload: "example.com",
			Chains: []string{"DIRECT", "香港 01"}, Download: 6 << 20, Upload: 800 << 10,
			DownloadSpeed: 240 << 10, UploadSpeed: 18 << 10,
			Start: time.Now().Add(-2 * time.Minute),
		},
		{
			ID: "c2", Host: "api.github.com", DestinationIP: "140.82.112.4", DestPort: 443,
			Process: "curl", Network: "tcp", Type: "HTTPS",
			SourceIP: "127.0.0.1", SourcePort: 52380, Rule: "DOMAIN", RulePayload: "api.github.com",
			Chains: []string{"香港 01"}, Download: 12 << 10, Upload: 400,
			Start: time.Now().Add(-30 * time.Second),
		},
	}
	m.Closed = []Connection{
		{
			ID: "c3", Host: "old.example", DestPort: 80, Process: "curl", Network: "tcp", Type: "HTTP",
			Rule: "MATCH", Chains: []string{"DIRECT"}, Download: 128, Upload: 64,
			Start: time.Now().Add(-time.Hour), Closed: true,
		},
	}
	return m
}

// CollapseSidebar is the collapsed navigation shared by the
// sidebar-collapsed screenshots and by the 36/56 brand checks.
// The rail is 72px, the logo area is 56×56, the mark stays 36×36,
// and the wordmark is hidden. A nil model starts from DemoModel.
func CollapseSidebar(m *Model) *Model {
	if m == nil {
		m = DemoModel()
	}
	m.Prefs.CollapseNav = true
	return m
}

// CollapsedBrandMeasure is what NewCollapsedTester shows for the brand slot.
type CollapsedBrandMeasure struct {
	SidebarW float32
	SlotW    float32
	SlotH    float32
	MarkW    float32
	MarkH    float32
	Wordmark bool
}

// MeasureCollapsedBrand reads the 72px rail, the 56×56 logo area, and the
// 36×36 mark from a frame made by NewCollapsedTester. Wordmark is true when
// 「Clash Go」 is still on screen.
func MeasureCollapsedBrand(tt *mygoui.Tester) (CollapsedBrandMeasure, bool) {
	var out CollapsedBrandMeasure
	if tt == nil {
		return out, false
	}
	side, okSide := tt.Find("sidebar")
	slot, okSlot := tt.Find("logo-slot")
	mark, okMark := tt.Find("mygo")
	if !okSide || !okSlot || !okMark {
		return out, false
	}
	out.SidebarW = side.W
	out.SlotW, out.SlotH = slot.W, slot.H
	out.MarkW, out.MarkH = mark.W, mark.H
	out.Wordmark = tt.HasText("Clash Go")
	return out, true
}

// NewCollapsedTester draws CollapseSidebar with ui.NewTester. No window is opened.
func NewCollapsedTester(m *Model, dark bool) *mygoui.Tester {
	return renderModel(CollapseSidebar(m), WindowWidth, WindowHeight, dark)
}

// NewPageTester draws m at the default window size with ui.NewTester.
// No window is opened. It is the same frame renderModel uses for screenshots.
func NewPageTester(m *Model) *mygoui.Tester {
	return renderModel(m, WindowWidth, WindowHeight, false)
}

// renderModel draws the model once the way the screenshot test does.
func renderModel(m *Model, width, height int, dark bool) *mygoui.Tester {
	if width <= 0 {
		width = WindowWidth
	}
	if height <= 0 {
		height = WindowHeight
	}
	tt := mygoui.NewTester(m.View, width, height)
	tt.SetScale(1)
	tt.SetDark(dark)
	for i := 0; i < 5 && tt.FramePending(); i++ {
		tt.Advance(time.Second)
	}
	return tt
}
