package ui

import (
	"regexp"
	"strings"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

// maxLogRows matches Verge's MAX_LOG_NUM.
const maxLogRows = 1000

// LogItem is one row on the logs page. Type is the kernel type string
// (debug, info, warning, error). Time is already formatted.
type LogItem struct {
	ID      uint64
	Time    string
	Type    string
	Payload string
}

// ClashLog is Verge's clash-log record: whether collection is on, the
// level menu value, and asc/desc order.
type ClashLog struct {
	Enable    bool
	LogFilter string
	LogOrder  string
}

// DefaultClashLog matches use-clash-log.ts: collecting, every level, oldest first.
func DefaultClashLog() ClashLog {
	return ClashLog{Enable: true, LogFilter: "all", LogOrder: "asc"}
}

// logFilterChoices are the logs.tsx menu values. Labels are
// shared.filters.logLevels.*. warn matches a type of warning, and err
// matches error, because the page filters with type.includes(logState).
var logFilterChoices = []struct{ Value, Label string }{
	{"all", "ALL"},
	{"debug", "DEBUG"},
	{"info", "INFO"},
	{"warn", "WARN"},
	{"err", "ERROR"},
}

func (m *Model) logsPage(c *mygoui.Context, pal Palette) {
	m.pollLogs(c)
	m.ensureLogIDs()
	rows := m.visibleLogs()
	w, _ := c.Size()
	mygoui.Column(c).Fill().MaxWidth(w - m.sidebarWidth()).Children(func() {
		m.logHeader(c, pal)
		m.logToolbar(c, pal)
		if m.Notice != "" {
			mygoui.Text(c, m.Notice).FontSize(13).TextColor(pal.Error).Padding(4, 20)
		}
		if logRegexInvalid(m.LogQuery, m.LogSearch) {
			mygoui.Text(c, "无效的正则表达式").FontSize(12).TextColor(pal.Error).Padding(0, 20, 4, 20)
		}
		if len(rows) == 0 {
			m.logEmpty(c, pal)
			return
		}
		m.logListView(c, pal, rows)
	})
}

func (m *Model) pollLogs(c *mygoui.Context) {
	if m.Hooks.RefreshLogs == nil {
		return
	}
	now := c.Now()
	if m.logSynced.IsZero() || now.Sub(m.logSynced) >= 200*time.Millisecond {
		m.Hooks.RefreshLogs()
		m.logSynced = now
	}
	c.After(200 * time.Millisecond)
}

func (m *Model) logHeader(c *mygoui.Context, pal Palette) {
	head := mygoui.Row(c).FillWidth().Height(58).Padding(0, 16).Gap(8).AlignItems(mygoui.Center).Shrink(0)
	head.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
	head.Children(func() {
		mygoui.Text(c, "日志").FontSize(PageTitleSize).Bold().Shrink(0).SingleLine()
		mygoui.Spacer(c).Grow(1)
		pause := "暂停"
		if !m.ClashLog.Enable {
			pause = "继续"
		}
		if textButton(c, pal, pause, pause, "log-pause") {
			m.ClashLog.Enable = !m.ClashLog.Enable
		}
		orderText := "顺序"
		orderTip := "按时间倒序"
		if m.ClashLog.LogOrder == "desc" {
			orderText = "倒序"
			orderTip = "按时间正序"
		}
		if textButton(c, pal, orderText, orderTip, "log-order") {
			m.toggleLogOrder()
		}
		if filledButton(c, pal, "清除", "log-clear", true) {
			m.clearLogs()
		}
	})
}

func (m *Model) logToolbar(c *mygoui.Context, pal Palette) {
	bar := mygoui.Row(c).FillWidth().Height(39).Padding(0, 10).Gap(8).AlignItems(mygoui.Center).Shrink(0).Label("log-toolbar")
	bar.Children(func() {
		cur := logFilterLabel(m.ClashLog.LogFilter)
		sel := mygoui.Select(c, &cur, logFilterLabels()).Width(128).Label("log-filter")
		if sel.Changed() {
			m.ClashLog.LogFilter = logFilterValue(cur)
		}
		mygoui.SearchField(c, &m.LogQuery).Grow(1).Label("log-search").Placeholder("过滤条件")
		m.connSearchToggle(c, pal, "大小写", "区分大小写", "log-case", &m.LogSearch.MatchCase)
		m.connSearchToggle(c, pal, "全字", "全字匹配", "log-word", &m.LogSearch.WholeWord)
		m.connSearchToggle(c, pal, "正则", "使用正则表达式", "log-regex", &m.LogSearch.Regex)
	})
}

func (m *Model) logEmpty(c *mygoui.Context, pal Palette) {
	mygoui.Column(c).Grow(1).FillWidth().AlignItems(mygoui.Center).Justify(mygoui.Center).Gap(8).Children(func() {
		mygoui.Column(c).Size(64, 48).Radius(8).Border(2, pal.TextMuted)
		mygoui.Text(c, "空空如也").FontSize(16).TextColor(pal.TextMuted)
	})
}

func (m *Model) logListView(c *mygoui.Context, pal Palette, rows []LogItem) {
	m.logList.Key = func(i int) any { return rows[i].ID }
	m.logList.Label = func(i int) string { return rows[i].Payload }
	// Ascending keeps the newest row in view while the viewport is at the
	// bottom, and glides there. Descending shows the reversed list from
	// the top; new lines land above the rows the user is reading.
	// Scrolling away from the end stops the follow. The list treats "at
	// the end" as about half a DIP, tighter than Verge's 20px.
	m.logList.FollowEnd = m.ClashLog.LogOrder != "desc"
	if m.logList.FollowEnd {
		m.logList.FollowAnimation = &mygoui.Smooth
	} else {
		m.logList.FollowAnimation = nil
	}
	list := mygoui.List(c, &m.logList, len(rows)).Grow(1).FillWidth()
	list.Rows(func(row mygoui.ListRow) {
		m.logRow(row.Context, pal, rows[row.Index])
	})
}

func (m *Model) logRow(c *mygoui.Context, pal Palette, item LogItem) {
	mygoui.Column(c).FillWidth().Padding(8, 12).Gap(2).BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider).Children(func() {
		mygoui.Row(c).Gap(8).AlignItems(mygoui.Center).Children(func() {
			mygoui.Text(c, item.Time).FontSize(14).TextColor(pal.TextMuted).SingleLine().Selectable()
			mygoui.Text(c, strings.ToUpper(item.Type)).FontSize(14).FontWeight(600).TextColor(logTone(pal, item.Type)).SingleLine().Selectable()
		})
		mygoui.Text(c, item.Payload).FontSize(14).TextColor(pal.Text).Selectable()
	})
}

func (m *Model) toggleLogOrder() {
	if m.ClashLog.LogOrder == "desc" {
		m.ClashLog.LogOrder = "asc"
		m.logList.ScrollToEnd()
		return
	}
	m.ClashLog.LogOrder = "desc"
	m.logList.ScrollTo(0, mygoui.Start)
}

func (m *Model) clearLogs() {
	m.Logs = nil
	if m.Hooks.ClearLogs != nil {
		m.Hooks.ClearLogs()
	}
}

// AppendLogs records lines while collection is enabled. Paused collection
// drops them. IDs are assigned here so the virtual list can follow a row.
func (m *Model) AppendLogs(items []LogItem) {
	if m == nil || !m.ClashLog.Enable || len(items) == 0 {
		return
	}
	for _, item := range items {
		m.logNext++
		item.ID = m.logNext
		m.Logs = append(m.Logs, item)
	}
	if extra := len(m.Logs) - maxLogRows; extra > 0 {
		m.Logs = append([]LogItem(nil), m.Logs[extra:]...)
	}
}

func (m *Model) ensureLogIDs() {
	for i := range m.Logs {
		if m.Logs[i].ID != 0 {
			if m.Logs[i].ID > m.logNext {
				m.logNext = m.Logs[i].ID
			}
			continue
		}
		m.logNext++
		m.Logs[i].ID = m.logNext
	}
}

func (m *Model) visibleLogs() []LogItem {
	filtered := filterLogs(m.Logs, m.ClashLog.LogFilter, m.LogQuery, m.LogSearch)
	if m.ClashLog.LogOrder != "desc" {
		return filtered
	}
	return reverseLogs(filtered)
}

func logTone(pal Palette, typ string) mygoui.Color {
	switch strings.ToLower(typ) {
	case "error", "err":
		return pal.Error
	case "warning", "warn":
		return pal.Warning
	case "info", "inf":
		return pal.Info
	default:
		return pal.Text
	}
}

func logFilterLabels() []string {
	out := make([]string, len(logFilterChoices))
	for i, choice := range logFilterChoices {
		out[i] = choice.Label
	}
	return out
}

func logFilterLabel(value string) string {
	for _, choice := range logFilterChoices {
		if choice.Value == value {
			return choice.Label
		}
	}
	return "ALL"
}

func logFilterValue(label string) string {
	for _, choice := range logFilterChoices {
		if choice.Label == label {
			return choice.Value
		}
	}
	return "all"
}

// filterLogs keeps rows whose type contains the menu value (or all of
// them) and whose lowercased "time type payload" passes the search box.
// An invalid regular expression matches nothing.
func filterLogs(items []LogItem, level, query string, opt SearchOpts) []LogItem {
	query = strings.TrimSpace(query)
	var re *regexp.Regexp
	if query != "" {
		compiled, ok := connRegexp(query, opt)
		if !ok {
			return nil
		}
		re = compiled
	}
	var out []LogItem
	for _, item := range items {
		if !logLevelMatch(item.Type, level) {
			continue
		}
		if query != "" && !nameMatch(item.searchText(), query, opt, re) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func logLevelMatch(typ, state string) bool {
	if state == "" || state == "all" {
		return true
	}
	return strings.Contains(typ, state)
}

func (item LogItem) searchText() string {
	return strings.ToLower(item.Time + " " + item.Type + " " + item.Payload)
}

func reverseLogs(items []LogItem) []LogItem {
	if len(items) == 0 {
		return items
	}
	out := make([]LogItem, len(items))
	for i := range items {
		out[len(items)-1-i] = items[i]
	}
	return out
}

func logRegexInvalid(query string, opt SearchOpts) bool {
	return connRegexInvalid(query, opt)
}
