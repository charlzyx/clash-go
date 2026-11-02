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
	return m
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
