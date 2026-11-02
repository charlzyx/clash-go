package ui

import (
	"strconv"
	"strings"

	"github.com/charlzyx/clash-go/internal/sitetest"
	"github.com/charlzyx/clash-go/internal/unlock"

	mygoui "github.com/egoist/mygo/ui"
)

// UnlockItem is one card on the unlock page.
// Delay -1 is untested, -2 is in progress, and >=0 is milliseconds.
type UnlockItem struct {
	Name   string
	Status string
	Region string
	Delay  int
}

// DefaultUnlocks is the 13-card Verge grid, all untested.
func DefaultUnlocks() []UnlockItem {
	names := unlock.Names()
	out := make([]UnlockItem, len(names))
	for i, name := range names {
		out[i] = UnlockItem{Name: name, Delay: -1}
	}
	return out
}

func (m *Model) unlockPage(c *mygoui.Context, pal Palette) {
	if m.Unlocks == nil {
		m.Unlocks = DefaultUnlocks()
	}
	busy := m.unlockBusy()
	m.page(c, pal, "解锁测试", func() {
		label := "测试全部"
		if busy {
			label = "测试中..."
		}
		if filledButton(c, pal, label, "unlock-all", !busy && len(m.Unlocks) > 0) && !busy {
			m.unlockAll()
		}
	}, func() {
		if len(m.Unlocks) == 0 {
			mygoui.Text(c, "暂无解锁测试项目").FontSize(14).TextColor(pal.TextMuted).Padding(24)
			return
		}
		mygoui.Grid(c).Columns(3).Gap(8).Padding(16).Children(func() {
			for i := range m.Unlocks {
				m.unlockCard(c, pal, i)
			}
		})
	})
}

func (m *Model) unlockBusy() bool {
	for _, it := range m.Unlocks {
		if it.Delay == -2 {
			return true
		}
	}
	return false
}

func (m *Model) unlockCard(c *mygoui.Context, pal Palette, i int) {
	it := m.Unlocks[i]
	b := mygoui.ButtonBase(c.Key("unlock-" + it.Name)).MinHeight(96).Radius(8).Padding(8).Label("unlock-" + it.Name)
	b.Background(pal.Card).BorderWidth(0, 0, 0, 4).BorderColor(unlockTone(pal, it))
	b.Children(func() {
		mygoui.Column(c).FillWidth().Gap(4).Children(func() {
			mygoui.Text(c, it.Name).FontSize(14).Bold().SingleLine()
			chip(c, pal, unlockChip(it), unlockToneName(it), false)
			region := "--"
			if it.Region != "" {
				region = flagEmoji(it.Region) + " " + it.Region
			}
			mygoui.Text(c, region).FontSize(12).TextColor(pal.TextMuted)
			mygoui.Text(c, unlockTime(it)).FontSize(12).TextColor(pal.TextMuted)
		})
	})
	if b.Clicked() {
		m.unlockOne(i)
	}
}

func (m *Model) unlockOne(i int) {
	if i < 0 || i >= len(m.Unlocks) || m.Unlocks[i].Delay == -2 {
		return
	}
	m.Unlocks[i].Delay = -2
	m.Unlocks[i].Status = "Pending"
	var status, region string
	var ms int
	var err error
	if m.Hooks.CheckUnlock != nil {
		status, region, ms, err = m.Hooks.CheckUnlock(m.Unlocks[i].Name)
	}
	if err != nil {
		m.Unlocks[i].Delay = -1
		m.Unlocks[i].Status = "Failed"
		if sitetest.IsTimeout(err) {
			m.Notice = "检测超时或失败"
			return
		}
		m.Notice = m.Unlocks[i].Name + " 检测失败"
		return
	}
	if status == "" {
		status = "Yes"
	}
	if ms < 0 {
		ms = 0
	}
	m.Unlocks[i].Status = status
	m.Unlocks[i].Region = region
	m.Unlocks[i].Delay = ms
}

func (m *Model) unlockAll() {
	if m.unlockBusy() || len(m.Unlocks) == 0 {
		return
	}
	for i := range m.Unlocks {
		m.unlockOne(i)
	}
}

func unlockChip(it UnlockItem) string {
	if it.Delay == -2 {
		return "测试中..."
	}
	switch it.Status {
	case "Yes":
		return "支持"
	case "No":
		return "不支持"
	case "Risky":
		return "存疑"
	case "Failed":
		return "测试失败"
	case "Soon":
		return "待检测"
	case "Completed":
		return "检测完成"
	default:
		return "待检测"
	}
}

func unlockTime(it UnlockItem) string {
	if it.Delay == -2 {
		return "测试中..."
	}
	if it.Delay >= 0 && it.Status != "" && it.Status != "Pending" {
		return strconv.Itoa(it.Delay) + " ms"
	}
	return "-- --"
}

func unlockToneName(it UnlockItem) Tone {
	if it.Delay == -2 || it.Status == "" || it.Status == "Pending" {
		return ToneMuted
	}
	switch it.Status {
	case "Yes":
		return ToneSuccess
	case "No", "Failed":
		return ToneError
	case "Risky", "Soon":
		return ToneWarning
	case "Completed":
		return ToneInfo
	default:
		return ToneMuted
	}
}

func unlockTone(pal Palette, it UnlockItem) mygoui.Color {
	return pal.Color(unlockToneName(it))
}

func flagEmoji(region string) string {
	region = strings.ToUpper(strings.TrimSpace(region))
	if len(region) != 2 {
		return ""
	}
	r := []rune(region)
	if r[0] < 'A' || r[0] > 'Z' || r[1] < 'A' || r[1] > 'Z' {
		return ""
	}
	return string(rune(r[0]-'A'+0x1F1E6)) + string(rune(r[1]-'A'+0x1F1E6))
}
