package ui

import (
	"fmt"
	"strings"

	mygoui "github.com/egoist/mygo/ui"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// UnlockItem is one media unlock row. The JSON shape matches Verge's
// UnlockItem so the prefs cache can stand in for
// localStorage['clash_verge_unlock_results'].
type UnlockItem struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Region    string `json:"region,omitempty"`
	CheckTime string `json:"check_time,omitempty"`
}

// UnlockStatusStyle is the chip, border, and icon for one status string.
// Border is success, error, warning, info, or divider.
// Icon is pending, yes, no, soon, or help.
type UnlockStatusStyle struct {
	Label  string
	Tone   Tone
	Border string
	Icon   string
	Bold   bool
}

// DefaultUnlockItems is getUnlockItems: the v2.5.9 catalog, every row Pending.
// The order is the checker list; the page sorts it by name before drawing.
func DefaultUnlockItems() []UnlockItem {
	names := []string{
		"哔哩哔哩大陆",
		"哔哩哔哩港澳台",
		"ChatGPT Web",
		"Claude",
		"Gemini",
		"Muse AI",
		"YouTube Premium",
		"Bahamut Anime",
		"Netflix",
		"Disney+",
		"Prime Video",
		"Spotify",
		"TikTok",
	}
	out := make([]UnlockItem, len(names))
	for i, name := range names {
		out[i] = UnlockItem{Name: name, Status: "Pending"}
	}
	return out
}

// UnlockColumns is the Verge grid: one column below 650, two below 900,
// three at the md breakpoint. The width is the window width, the same
// breakpoint home uses.
func UnlockColumns(windowWidth float32) int {
	switch {
	case windowWidth >= BreakpointMD:
		return 3
	case windowWidth >= BreakpointSM:
		return 2
	default:
		return 1
	}
}

func normalizeUnlockName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func unlockStatusPriority(status string) int {
	if status == "Pending" {
		return 0
	}
	return 1
}

// mergeUnlockFields keeps preferred and fills a missing region or time
// from fallback. Empty means the JSON field was absent.
func mergeUnlockFields(preferred, fallback UnlockItem) UnlockItem {
	if preferred.Region == "" {
		preferred.Region = fallback.Region
	}
	if preferred.CheckTime == "" {
		preferred.CheckTime = fallback.CheckTime
	}
	return preferred
}

// DedupeUnlockItems merges same-name rows. Comparison ignores case and
// surrounding space. A non-Pending row overwrites Pending. Two rows with
// the same priority keep the later one, and a missing region or time is
// filled from the one that had it.
func DedupeUnlockItems(items []UnlockItem) []UnlockItem {
	order := make([]string, 0, len(items))
	got := make(map[string]UnlockItem, len(items))
	for _, item := range items {
		key := normalizeUnlockName(item.Name)
		existing, ok := got[key]
		if !ok {
			got[key] = item
			order = append(order, key)
			continue
		}
		switch {
		case unlockStatusPriority(item.Status) > unlockStatusPriority(existing.Status):
			got[key] = mergeUnlockFields(item, existing)
		case unlockStatusPriority(item.Status) < unlockStatusPriority(existing.Status):
			got[key] = mergeUnlockFields(existing, item)
		default:
			got[key] = mergeUnlockFields(item, existing)
		}
	}
	out := make([]UnlockItem, 0, len(order))
	for _, key := range order {
		out = append(out, got[key])
	}
	return out
}

// MergeUnlockItems applies a cached list onto the catalog. Stored names
// that the catalog does not have are dropped. The catalog name wins so
// "netflix" in the cache still displays "Netflix". The result is sorted
// by name with the Chinese collator, the stand-in for localeCompare.
func MergeUnlockItems(defaults, stored []UnlockItem) []UnlockItem {
	stored = DedupeUnlockItems(stored)
	existing := make(map[string]UnlockItem, len(stored))
	for _, item := range stored {
		existing[normalizeUnlockName(item.Name)] = item
	}
	merged := make([]UnlockItem, 0, len(defaults))
	for _, item := range defaults {
		if matched, ok := existing[normalizeUnlockName(item.Name)]; ok {
			matched.Name = item.Name
			merged = append(merged, matched)
			continue
		}
		merged = append(merged, item)
	}
	return SortUnlockItems(merged)
}

// SortUnlockItems orders a copy by name. It does not dedupe.
func SortUnlockItems(items []UnlockItem) []UnlockItem {
	out := append([]UnlockItem(nil), items...)
	c := collate.New(language.Chinese)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && c.CompareString(out[j].Name, out[j-1].Name) < 0 {
			out[j], out[j-1] = out[j-1], out[j]
			j--
		}
	}
	return out
}

// UnlockStatusStyleOf maps a raw status onto the chip table in unlock.tsx.
// Unknown strings keep their own text, the default chip, a divider border,
// and the help icon. "Failed" anywhere in the string is the failure color
// before the ISP / region special cases.
func UnlockStatusStyleOf(status string) UnlockStatusStyle {
	st := UnlockStatusStyle{
		Label:  unlockStatusLabel(status),
		Tone:   ToneMuted,
		Border: "divider",
		Icon:   "help",
		Bold:   status != "Pending",
	}
	switch {
	case status == "Pending":
		st.Icon = "pending"
		st.Bold = false
	case status == "Yes":
		st.Tone, st.Border, st.Icon = ToneSuccess, "success", "yes"
	case status == "No":
		st.Tone, st.Border, st.Icon = ToneError, "error", "no"
	case status == "Soon":
		st.Tone, st.Border, st.Icon = ToneWarning, "warning", "soon"
	case status == "Risky":
		st.Tone, st.Border, st.Icon = ToneWarning, "warning", "help"
	case strings.Contains(status, "Failed"):
		st.Tone, st.Border, st.Icon = ToneError, "error", "help"
	case status == "Completed":
		st.Tone, st.Border, st.Icon = ToneInfo, "info", "help"
	case status == "Disallowed ISP" || status == "Blocked" || status == "Unsupported Country/Region":
		st.Tone, st.Border, st.Icon = ToneError, "divider", "help"
	}
	return st
}

func unlockStatusLabel(status string) string {
	switch status {
	case "Pending":
		return "待检测"
	case "Yes":
		return "支持"
	case "No":
		return "不支持"
	case "Failed":
		return "测试失败"
	case "Risky":
		return "存疑"
	case "Completed":
		return "检测完成"
	case "Disallowed ISP":
		return "不允许的 ISP"
	case "Originals Only":
		return "仅限原创"
	case "No (IP Banned By Disney+)":
		return "不支持（IP 被 Disney+ 禁止）"
	case "Unsupported Country/Region":
		return "不支持的国家/地区"
	case "Failed (Network Connection)":
		return "测试失败（网络连接问题）"
	default:
		return status
	}
}

func unlockBorderColor(pal Palette, border string) mygoui.Color {
	switch border {
	case "success":
		return pal.Success
	case "error":
		return pal.Error
	case "warning":
		return pal.Warning
	case "info":
		return pal.Info
	default:
		return pal.Divider
	}
}

// syncUnlock loads the catalog when the route is entered. Leaving and
// coming back loads again, the way the page remounts. Hiding the window
// while the route stays put does not load again; Verge's effect has an
// empty dependency list and does not follow visibility.
func (m *Model) syncUnlock() {
	if m == nil || m.Router == nil {
		return
	}
	on := m.Router.Path() == "/unlock"
	if on && !m.unlockOnPage {
		m.loadUnlock()
	}
	m.unlockOnPage = on
}

func (m *Model) loadUnlock() {
	if m.Hooks.RefreshUnlock != nil {
		m.Hooks.RefreshUnlock()
		return
	}
	m.Unlock = MergeUnlockItems(DefaultUnlockItems(), m.Prefs.UnlockResults)
}

func (m *Model) unlockPage(c *mygoui.Context, pal Palette) {
	w, _ := c.Size()
	mygoui.Column(c).Fill().Children(func() {
		head := mygoui.Row(c).Height(58).Padding(0, 20).Gap(8).AlignItems(mygoui.Center).Shrink(0)
		head.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
		head.Children(func() {
			mygoui.Text(c, "解锁测试").FontSize(PageTitleSize).Bold().Shrink(0)
			mygoui.Spacer(c).Grow(1)
			m.unlockTestAllButton(c, pal)
		})
		if len(m.Unlock) == 0 {
			m.unlockEmpty(c, pal)
			return
		}
		mygoui.Scroll(c).Grow(1).Children(func() {
			cols := UnlockColumns(w)
			mygoui.Grid(c).Columns(cols).Gap(12).Padding(12).Label("unlock-grid").Children(func() {
				for i := range m.Unlock {
					m.unlockCard(c, pal, m.Unlock[i])
				}
			})
		})
	})
}

func (m *Model) unlockTestAllButton(c *mygoui.Context, pal Palette) {
	busy := m.unlockAll || len(m.unlockLoading) > 0
	disabled := len(m.Unlock) == 0 || busy
	b := mygoui.ButtonBase(c.Key("unlock-test-all")).Height(30).Padding(0, 12).Gap(6).Radius(6).
		Disabled(disabled).Label("unlock-test-all").AlignItems(mygoui.Center)
	if disabled && !m.unlockAll {
		b.Background(pal.Primary.Alpha(0.38)).TextColor(pal.OnPrimary.Alpha(0.7))
	} else {
		b.Background(pal.Primary).TextColor(pal.OnPrimary)
	}
	b.Children(func() {
		if m.unlockAll {
			mygoui.Spinner(c).Size(16, 16).TextColor(pal.OnPrimary)
			mygoui.Text(c, m.unlockTestingLabel()).FontSize(13)
			return
		}
		mygoui.Icon(c, unlockIconRefresh).Size(16, 16)
		mygoui.Text(c, "全部测试").FontSize(13)
	})
	if b.Clicked() {
		m.testAllUnlock()
	}
}

func (m *Model) unlockTestingLabel() string {
	n := len(m.Unlock)
	if n == 0 {
		return "测试中..."
	}
	done := n - len(m.unlockLoading)
	if done < 0 {
		done = 0
	}
	if done > n {
		done = n
	}
	return fmt.Sprintf("测试中... %d/%d", done, n)
}

func (m *Model) unlockEmpty(c *mygoui.Context, pal Palette) {
	// 连接页同一套空态：方框代替 InboxRounded，文案是 BaseEmpty 的默认「空空如也」。
	mygoui.Column(c).Grow(1).FillWidth().AlignItems(mygoui.Center).Justify(mygoui.Center).Gap(8).
		Label("unlock-empty").Children(func() {
		mygoui.Column(c).Size(64, 48).Radius(8).Border(2, pal.TextMuted)
		mygoui.Text(c, "空空如也").FontSize(16).TextColor(pal.TextMuted)
	})
}

func (m *Model) unlockCard(c *mygoui.Context, pal Palette, item UnlockItem) {
	st := UnlockStatusStyleOf(item.Status)
	loading := m.unlockLoading[item.Name]
	disabled := m.unlockAll || loading
	mygoui.Column(c).FillWidth().MinHeight(108).Radius(8).Background(pal.Card).Clip().
		Border(1, pal.Divider).Label("unlock-card-" + item.Name).Children(func() {
		mygoui.Row(c).FillWidth().Grow(1).Children(func() {
			mygoui.Box(c).Width(4).FillHeight().Shrink(0).Background(unlockBorderColor(pal, st.Border))
			b := mygoui.ButtonBase(c.Key("unlock-"+item.Name)).Grow(1).Padding(10, 12).
				Disabled(disabled).Label("unlock-" + item.Name).Tooltip("测试")
			b.Background(mygoui.Transparent).AlignItems(mygoui.Start)
			b.Children(func() {
				mygoui.Column(c).Grow(1).Gap(8).FillWidth().Children(func() {
					mygoui.Row(c).FillWidth().Gap(8).AlignItems(mygoui.Center).Children(func() {
						mygoui.Text(c, item.Name).FontSize(16).FontWeight(600).Grow(1).SingleLine()
						if loading {
							mygoui.Spinner(c).Size(16, 16)
						} else {
							mygoui.Icon(c, unlockIconRefresh).Size(18, 18).TextColor(pal.Primary)
						}
					})
					mygoui.Row(c).Gap(8).AlignItems(mygoui.Center).Children(func() {
						unlockStatusChip(c, pal, st)
						if item.Region != "" {
							mygoui.Row(c).Padding(2, 8).Radius(16).Shrink(0).Border(1, pal.Info).Children(func() {
								mygoui.Text(c, item.Region).FontSize(12).TextColor(pal.Info)
							})
						}
					})
				})
			})
			if b.Clicked() {
				m.testOneUnlock(item.Name)
			}
		})
		mygoui.Box(c).Height(1).FillWidth().Margin(0, 8).Background(pal.Divider.Alpha(0.5))
		when := item.CheckTime
		if when == "" {
			when = "-- --"
		}
		mygoui.Row(c).Padding(2, 12, 6, 12).Children(func() {
			mygoui.Spacer(c).Grow(1)
			mygoui.Text(c, when).FontSize(11).TextColor(pal.TextMuted)
		})
	})
}

func unlockStatusChip(c *mygoui.Context, pal Palette, st UnlockStatusStyle) {
	fg := pal.Color(st.Tone)
	bg := fg.Alpha(0.12)
	if st.Icon == "pending" {
		fg = pal.Text
		bg = pal.Text.Alpha(0.08)
	}
	weight := 400
	if st.Bold {
		weight = 700
	}
	mygoui.Row(c).Padding(2, 8).Gap(4).Radius(16).Shrink(0).Background(bg).
		AlignItems(mygoui.Center).Children(func() {
		mygoui.Icon(c, unlockStatusIcon(st.Icon)).Size(14, 14).TextColor(fg)
		mygoui.Text(c, st.Label).FontSize(12).FontWeight(weight).TextColor(fg)
	})
}

func (m *Model) testAllUnlock() {
	if m == nil || len(m.Unlock) == 0 || m.unlockAll || len(m.unlockLoading) > 0 {
		return
	}
	m.SetNotice("")
	m.unlockAll = true
	if m.unlockLoading == nil {
		m.unlockLoading = map[string]bool{}
	}
	for _, item := range m.Unlock {
		m.unlockLoading[item.Name] = true
	}
	if m.Hooks.CheckUnlock == nil {
		m.SetNotice("检测超时或失败")
		m.finishUnlockAll(false)
		return
	}
	items, err := m.Hooks.CheckUnlock("", func(item UnlockItem) {
		m.replaceUnlock(item)
		delete(m.unlockLoading, item.Name)
	})
	if err != nil {
		m.SetNotice("检测超时或失败")
		m.finishUnlockAll(false)
		return
	}
	if items != nil {
		m.Unlock = SortUnlockItems(DedupeUnlockItems(items))
	} else {
		m.Unlock = SortUnlockItems(DedupeUnlockItems(m.Unlock))
	}
	m.finishUnlockAll(true)
}

func (m *Model) finishUnlockAll(save bool) {
	m.unlockLoading = map[string]bool{}
	m.unlockAll = false
	if save {
		m.persistUnlock()
	}
}

func (m *Model) testOneUnlock(name string) {
	if m == nil || name == "" || m.unlockAll || m.unlockLoading[name] {
		return
	}
	if m.unlockLoading == nil {
		m.unlockLoading = map[string]bool{}
	}
	m.unlockLoading[name] = true
	m.SetNotice("")
	if m.Hooks.CheckUnlock == nil {
		m.SetNotice(name + " 检测失败")
		delete(m.unlockLoading, name)
		return
	}
	items, err := m.Hooks.CheckUnlock(name, func(item UnlockItem) {
		m.replaceUnlock(item)
	})
	if err != nil {
		m.SetNotice(name + " 检测失败")
		delete(m.unlockLoading, name)
		return
	}
	if len(items) > 0 {
		for _, item := range items {
			if item.Name == "" {
				item.Name = name
			}
			if item.Name == name {
				m.replaceUnlock(item)
			}
		}
	}
	delete(m.unlockLoading, name)
	m.persistUnlock()
}

func (m *Model) replaceUnlock(item UnlockItem) {
	for i := range m.Unlock {
		if m.Unlock[i].Name == item.Name {
			m.Unlock[i] = item
			return
		}
	}
}

func (m *Model) persistUnlock() {
	m.Prefs.UnlockResults = DedupeUnlockItems(append([]UnlockItem(nil), m.Unlock...))
	m.savePrefs()
}

const unlockIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">`

var (
	unlockIconPending = mustNav(unlockIconSVG + `<circle cx="12" cy="12" r="8"/><path d="M12 8 V12 L15 14"/></svg>`)
	unlockIconYes     = mustNav(unlockIconSVG + `<circle cx="12" cy="12" r="8"/><path d="M8.5 12.2 L11 14.7 L15.5 9.5"/></svg>`)
	unlockIconNo      = mustNav(unlockIconSVG + `<circle cx="12" cy="12" r="8"/><path d="M9 9 L15 15"/><path d="M15 9 L9 15"/></svg>`)
	unlockIconSoon    = mustNav(unlockIconSVG + `<circle cx="12" cy="12" r="8"/><path d="M12 8.5 V12.2 H15.2"/></svg>`)
	unlockIconHelp    = mustNav(unlockIconSVG + `<circle cx="12" cy="12" r="8"/><path d="M9.2 9.6 C9.6 8 14.2 7.6 14.4 10.2 C14.6 12.2 12 12.2 12 14.2"/><path d="M12 17 H12.1"/></svg>`)
	unlockIconRefresh = mustNav(unlockIconSVG + `<path d="M19.5 12 A7.5 7.5 0 1 1 16.8 6.2"/><path d="M19.5 4.5 V8.5 H15.5"/></svg>`)
)

func unlockStatusIcon(kind string) *mygoui.SVG {
	switch kind {
	case "pending":
		return unlockIconPending
	case "yes":
		return unlockIconYes
	case "no":
		return unlockIconNo
	case "soon":
		return unlockIconSoon
	default:
		return unlockIconHelp
	}
}
