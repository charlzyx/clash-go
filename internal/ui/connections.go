package ui

import (
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

// Connection is one proxied flow. Closed is true once it has left the
// core snapshot and is kept only for the 已关闭 tab.
type Connection struct {
	ID            string
	Host          string
	DestinationIP string
	DestPort      uint16
	RemoteDest    string
	Process       string
	ProcessPath   string
	Network       string
	Type          string
	SourceIP      string
	SourcePort    uint16
	Rule          string
	RulePayload   string
	Chains        []string
	Upload        int64
	Download      int64
	UploadSpeed   int64
	DownloadSpeed int64
	Start         time.Time
	Closed        bool
	ClosedAt      time.Time
}

type connColumn struct {
	ID    string
	Title string
	Width float32
	Right bool
}

// connectionColumns is the Verge table, in its default order. Every
// column starts visible. Widths are the Verge defaults; drag-resize is
// left to the table widget's own header and is not persisted.
var connectionColumns = []connColumn{
	{ID: "host", Title: "主机", Width: 180},
	{ID: "download", Title: "下载量", Width: 88, Right: true},
	{ID: "upload", Title: "上传量", Width: 88, Right: true},
	{ID: "dlSpeed", Title: "下载速度", Width: 100, Right: true},
	{ID: "ulSpeed", Title: "上传速度", Width: 100, Right: true},
	{ID: "chains", Title: "链路", Width: 220},
	{ID: "rule", Title: "规则", Width: 180},
	{ID: "process", Title: "进程", Width: 140},
	{ID: "time", Title: "连接时间", Width: 100, Right: true},
	{ID: "source", Title: "源地址", Width: 150},
	{ID: "destination", Title: "目标地址", Width: 150},
	{ID: "type", Title: "类型", Width: 120},
}

const connScrollTop = float32(100)

func (m *Model) connectionsPage(c *mygoui.Context, pal Palette) {
	m.pollConnections(c)
	now := c.Now()
	rows := m.visibleConnections()
	w, _ := c.Size()
	mygoui.Column(c).Fill().MaxWidth(w - m.sidebarWidth()).Children(func() {
		m.connHeader(c, pal, now)
		m.connToolbar(c, pal, rows)
		if m.Notice != "" {
			mygoui.Text(c, m.Notice).FontSize(13).TextColor(pal.Error).Padding(4, 20)
		}
		if connRegexInvalid(m.ConnQuery, m.ConnSearch) {
			mygoui.Text(c, "无效的正则表达式").FontSize(12).TextColor(pal.Error).Padding(0, 20, 4, 20)
		}
		if len(rows) == 0 {
			m.connEmpty(c, pal)
		} else if m.ConnLayout == "table" {
			m.connTableView(c, pal, rows, now)
		} else {
			m.connListView(c, pal, rows, now)
		}
		if m.connScroll.Y > connScrollTop {
			m.connTopButton(c, pal)
		}
	})
	m.connDetail(c, pal, now)
}

func (m *Model) pollConnections(c *mygoui.Context) {
	if m.Hooks.RefreshConnections == nil {
		return
	}
	now := c.Now()
	if m.connSynced.IsZero() || now.Sub(m.connSynced) >= time.Second {
		before := m.connSynced
		m.Hooks.RefreshConnections(now)
		if m.connSynced.Equal(before) {
			m.connSynced = now
		}
	}
	c.After(time.Second)
}

func (m *Model) connHeader(c *mygoui.Context, pal Palette, now time.Time) {
	head := mygoui.Row(c).FillWidth().Height(58).Padding(0, 16).Gap(8).AlignItems(mygoui.Center).Shrink(0)
	head.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
	head.Children(func() {
		mygoui.Text(c, "连接").FontSize(PageTitleSize).Bold().Shrink(0).SingleLine()
		mygoui.Spacer(c).Grow(1)
		mygoui.Text(c, "下载量 "+FormatTraffic(m.ConnDown)).FontSize(13).TextColor(pal.TextMuted).SingleLine().Shrink(0)
		mygoui.Text(c, "上传量 "+FormatTraffic(m.ConnUp)).FontSize(13).TextColor(pal.TextMuted).SingleLine().Shrink(0)
		if filledButton(c, pal, "列表", "conn-layout-list", m.ConnLayout != "table") {
			m.setConnLayout("list")
		}
		if filledButton(c, pal, "表格", "conn-layout-table", m.ConnLayout == "table") {
			m.setConnLayout("table")
		}
		if textButton(c, pal, "关闭全部", "关闭全部连接", "close-all") {
			m.closeAllConnections(now)
		}
	})
}

func (m *Model) connToolbar(c *mygoui.Context, pal Palette, rows []Connection) {
	bar := mygoui.Row(c).FillWidth().MinHeight(36).Padding(4, 12).Gap(8).AlignItems(mygoui.Center).Shrink(0)
	bar.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
	bar.Children(func() {
		if m.connTabButton(c, pal, "活跃", len(m.Active), m.connTab == 0, "conn-active") {
			m.setConnTab(0)
		}
		if m.connTabButton(c, pal, "已关闭", len(m.Closed), m.connTab == 1, "conn-closed") {
			m.setConnTab(1)
		}
		if m.connTab == 1 && textButton(c, pal, "清除", "清除已关闭连接", "conn-clear") {
			m.clearClosed()
		}
		if m.ConnLayout != "table" {
			cur := connOrderLabel(m.ConnOrder)
			sel := mygoui.Select(c, &cur, []string{"默认", "上传速度", "下载速度"}).Width(128).Label("conn-order")
			if sel.Changed() {
				m.ConnOrder = connOrderID(cur)
			}
		}
		mygoui.SearchField(c, &m.ConnQuery).Width(180).Label("conn-search").Placeholder("过滤条件")
		m.connSearchToggle(c, pal, "大小写", "区分大小写", "conn-case", &m.ConnSearch.MatchCase)
		m.connSearchToggle(c, pal, "全字", "全字匹配", "conn-word", &m.ConnSearch.WholeWord)
		m.connSearchToggle(c, pal, "正则", "使用正则表达式", "conn-regex", &m.ConnSearch.Regex)
		if m.ConnLayout == "table" && len(rows) > 0 {
			m.connColumnButton(c, pal)
		}
	})
}

func (m *Model) connTabButton(c *mygoui.Context, pal Palette, word string, n int, on bool, id string) bool {
	b := mygoui.ButtonBase(c.Key(id)).Height(28).Padding(0, 10).Gap(4).Radius(6).Label(id)
	if on {
		b.Background(pal.Primary).TextColor(pal.OnPrimary)
	} else {
		b.Background(mygoui.Transparent).TextColor(pal.Text).Border(1, pal.Divider)
	}
	b.AlignItems(mygoui.Center)
	b.Children(func() {
		mygoui.Text(c, word).FontSize(13)
		mygoui.Text(c, itoa(n)).FontSize(13)
	})
	return b.Clicked()
}

func (m *Model) connSearchToggle(c *mygoui.Context, pal Palette, label, tip, id string, on *bool) {
	b := mygoui.ButtonBase(c.Key(id)).Height(28).Padding(0, 8).Radius(6).Label(id).Tooltip(tip)
	if on != nil && *on {
		b.Background(pal.Primary).TextColor(pal.OnPrimary)
	} else {
		b.Background(mygoui.Transparent).TextColor(pal.Text).Border(1, pal.Divider)
	}
	b.AlignItems(mygoui.Center)
	b.Children(func() { mygoui.Text(c, label).FontSize(13) })
	if b.Clicked() && on != nil {
		*on = !*on
	}
}

func (m *Model) connColumnButton(c *mygoui.Context, pal Palette) {
	b := mygoui.ButtonBase(c.Key("conn-columns")).Height(28).Padding(0, 8).Radius(6).Label("列设置")
	if m.connColsOpen {
		b.Background(pal.Primary).TextColor(pal.OnPrimary)
	} else {
		b.Background(mygoui.Transparent).TextColor(pal.Text).Border(1, pal.Divider)
	}
	b.Children(func() { mygoui.Text(c, "列").FontSize(13) })
	if b.Clicked() {
		m.connColsOpen = true
	}
	mygoui.Popover(c, b, &m.connColsOpen, func() {
		mygoui.Column(c).Width(180).Gap(4).Padding(4).Children(func() {
			mygoui.Text(c, "列设置").FontSize(14).Bold()
			for _, col := range connectionColumns {
				shown := !m.connHidden[col.ID]
				box := mygoui.Checkbox(c, &shown, col.Title).Label("col-" + col.ID)
				if box.Changed() {
					m.setColumnVisible(col.ID, shown)
				}
			}
		})
	})
}

func (m *Model) connEmpty(c *mygoui.Context, pal Palette) {
	mygoui.Column(c).Grow(1).FillWidth().AlignItems(mygoui.Center).Justify(mygoui.Center).Gap(8).Children(func() {
		mygoui.Column(c).Size(64, 48).Radius(8).Border(2, pal.TextMuted)
		mygoui.Text(c, "空空如也").FontSize(16).TextColor(pal.TextMuted)
	})
}

func (m *Model) connListView(c *mygoui.Context, pal Palette, rows []Connection, now time.Time) {
	m.connList.Key = func(i int) any { return rows[i].ID }
	m.connList.KeepOffset = true
	m.connList.Label = func(i int) string { return rows[i].Endpoint() }
	list := mygoui.List(c, &m.connList, len(rows)).Grow(1).FillWidth().TrackScroll(&m.connScroll)
	list.Rows(func(row mygoui.ListRow) {
		item := rows[row.Index]
		m.connRow(row.Context, pal, item, now)
	})
}

func (m *Model) connRow(c *mygoui.Context, pal Palette, item Connection, now time.Time) {
	mygoui.Row(c).FillWidth().MinHeight(64).Padding(6, 12).Gap(8).AlignItems(mygoui.Center).
		BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider).Children(func() {
		body := mygoui.ButtonBase(c.Key("conn-row-" + item.ID)).Grow(1).Label(item.Endpoint())
		body.Background(mygoui.Transparent).AlignItems(mygoui.Center)
		body.Children(func() {
			mygoui.Column(c).Grow(1).Gap(2).Children(func() {
				mygoui.Text(c, item.Endpoint()).FontSize(14).SingleLine()
				mygoui.Row(c).Gap(4).Children(func() {
					for _, tag := range item.Tags() {
						chip(c, pal, tag, ToneMuted, false)
					}
				})
				mygoui.Row(c).Gap(8).Children(func() {
					mygoui.Text(c, "下载 "+FormatTraffic(item.Download)).FontSize(12).TextColor(pal.TextMuted)
					mygoui.Text(c, "上传 "+FormatTraffic(item.Upload)).FontSize(12).TextColor(pal.TextMuted)
					if item.DownloadSpeed >= 100 || item.UploadSpeed >= 100 {
						mygoui.Text(c, FormatSpeed(item.UploadSpeed)+" / "+FormatSpeed(item.DownloadSpeed)).FontSize(12).TextColor(pal.Primary)
					}
				})
			})
		})
		if body.Clicked() {
			m.openConnection(item)
		}
		mygoui.Text(c, RelativeSince(now, item.Start)).FontSize(12).TextColor(pal.TextMuted).Shrink(0)
		if !item.Closed {
			b := mygoui.ButtonBase(c.Key("conn-close-"+item.ID)).Size(28, 28).Radius(6).Label("conn-close-" + item.ID)
			b.Background(mygoui.Transparent).TextColor(pal.Error).Tooltip("关闭连接")
			b.AlignItems(mygoui.Center).Justify(mygoui.Center)
			b.Children(func() { mygoui.Text(c, "×").FontSize(16) })
			if b.Clicked() {
				m.closeOne(item.ID, now)
			}
		}
	})
}

func (m *Model) connTableView(c *mygoui.Context, pal Palette, rows []Connection, now time.Time) {
	cols := m.visibleColumns()
	spec := make([]mygoui.TableColumn, len(cols))
	for i, col := range cols {
		spec[i] = mygoui.TableColumn{
			Title: col.Title, ID: col.ID, Width: col.Width, MinWidth: 64, Sortable: true,
		}
		if col.Right {
			spec[i].Align = mygoui.End
		}
	}
	sortNow := m.connSort
	rows = sortConnections(rows, "", true, sortNow)
	m.connTable.Key = func(i int) any { return rows[i].ID }
	m.connTable.KeepOffset = true
	m.connTable.Selected = &m.connSelected
	m.connTable.Sort = &m.connSort
	m.connTable.Label = func(i int) string { return rows[i].Endpoint() }
	table := mygoui.Table(c, &m.connTable, spec, len(rows), func(row, col int) {
		id := spec[col].ID
		mygoui.Text(c, rows[row].columnText(id, now)).FontSize(13).SingleLine()
	}).Grow(1).FillWidth().Shrink(1).TrackScroll(&m.connScroll)
	if table.Changed() && m.connSelected >= 0 && m.connSelected < len(rows) {
		m.openConnection(rows[m.connSelected])
	}
}

func (m *Model) connTopButton(c *mygoui.Context, pal Palette) {
	b := mygoui.ButtonBase(c.Key("conn-top")).Absolute().Right(20).Bottom(20).Size(40, 40).Radius(20).Label("回到顶部")
	b.Background(pal.Primary).TextColor(pal.OnPrimary)
	b.AlignItems(mygoui.Center).Justify(mygoui.Center).Tooltip("回到顶部")
	b.Children(func() { mygoui.Text(c, "↑").FontSize(18).TextColor(pal.OnPrimary) })
	if b.Clicked() {
		m.connScroll.Y = 0
		m.connList.ScrollTo(0, mygoui.Start)
		m.connTable.ScrollTo(0, mygoui.Start)
	}
}

func (m *Model) connDetail(c *mygoui.Context, pal Palette, now time.Time) {
	item, ok := m.connectionByID(m.connDetailID)
	if !ok {
		return
	}
	mygoui.Overlay(c, func() {
		dim := mygoui.ButtonBase(c.Key("conn-dim")).Absolute().Top(0).Left(0).Right(0).Bottom(0).Label("conn-dim")
		dim.Background(mygoui.RGBA(0, 0, 0, 0.32))
		if dim.Clicked() {
			m.connDetailID = ""
		}
		panel := mygoui.Column(c).Absolute().Top(0).Bottom(0).Right(0).Width(380).Padding(16).Gap(8)
		panel.Background(pal.Paper).BorderWidth(0, 0, 0, 1).BorderColor(pal.Divider)
		panel.Children(func() {
			mygoui.Row(c).AlignItems(mygoui.Center).Gap(8).Children(func() {
				mygoui.Text(c, "连接详情").FontSize(18).Bold().Grow(1).SingleLine()
				if item.Closed {
					chip(c, pal, "已关闭", ToneMuted, false)
				} else {
					chip(c, pal, "活跃", ToneSuccess, true)
				}
				if textButton(c, pal, "关闭", "关闭详情", "conn-detail-close") {
					m.connDetailID = ""
				}
			})
			mygoui.Scroll(c).Grow(1).Children(func() {
				mygoui.Column(c).Gap(8).Padding(0, 0, 12, 0).Children(func() {
					for _, field := range item.detailFields(now) {
						mygoui.Column(c).Gap(2).Children(func() {
							mygoui.Text(c, field.label).FontSize(12).TextColor(pal.TextMuted)
							value := field.value
							if value == "" {
								value = "-"
							}
							mygoui.Text(c, value).FontSize(14)
						})
					}
				})
			})
			if !item.Closed {
				if filledButton(c, pal, "关闭连接", "conn-detail-delete", true) {
					m.closeOne(item.ID, now)
				}
			}
		})
	})
}

func (m *Model) visibleConnections() []Connection {
	src := m.Active
	if m.connTab == 1 {
		src = m.Closed
	}
	matched := filterConnections(src, m.ConnQuery, m.ConnSearch)
	if m.ConnLayout == "table" {
		return matched
	}
	return sortConnections(matched, m.ConnOrder, false, mygoui.SortOrder{})
}

func (m *Model) visibleColumns() []connColumn {
	var out []connColumn
	for _, col := range connectionColumns {
		if m.connHidden[col.ID] {
			continue
		}
		out = append(out, col)
	}
	if len(out) == 0 {
		return connectionColumns[:1]
	}
	return out
}

func (m *Model) setColumnVisible(id string, on bool) {
	if m.connHidden == nil {
		m.connHidden = map[string]bool{}
	}
	if on {
		delete(m.connHidden, id)
		return
	}
	visible := 0
	for _, col := range connectionColumns {
		if col.ID == id {
			continue
		}
		if !m.connHidden[col.ID] {
			visible++
		}
	}
	if visible == 0 {
		return
	}
	m.connHidden[id] = true
}

func (m *Model) setConnLayout(layout string) {
	if layout != "table" {
		layout = "list"
	}
	if m.ConnLayout == layout {
		return
	}
	m.ConnLayout = layout
	m.connColsOpen = false
	m.connSelected = -1
	m.connScroll.Y = 0
}

func (m *Model) setConnTab(tab int) {
	if tab != 1 {
		tab = 0
	}
	if m.connTab == tab {
		return
	}
	m.connTab = tab
	m.connDetailID = ""
	m.connColsOpen = false
	m.connSelected = -1
	m.connScroll.Y = 0
}

func (m *Model) openConnection(item Connection) {
	m.connDetailID = item.ID
}

func (m *Model) connectionByID(id string) (Connection, bool) {
	if id == "" {
		return Connection{}, false
	}
	for _, item := range m.Active {
		if item.ID == id {
			item.Closed = false
			return item, true
		}
	}
	for _, item := range m.Closed {
		if item.ID == id {
			item.Closed = true
			return item, true
		}
	}
	return Connection{}, false
}

func (m *Model) closeAllConnections(now time.Time) {
	if m.Hooks.CloseAllConnections != nil {
		if err := m.Hooks.CloseAllConnections(); err != nil {
			return
		}
	}
	m.markAllClosed(now)
}

func (m *Model) closeOne(id string, now time.Time) {
	if id == "" {
		return
	}
	if m.Hooks.CloseConnection != nil {
		if err := m.Hooks.CloseConnection(id); err != nil {
			return
		}
	}
	kept := make([]Connection, 0, len(m.Active))
	for _, item := range m.Active {
		if item.ID != id {
			kept = append(kept, item)
			continue
		}
		item.Closed = true
		if item.ClosedAt.IsZero() {
			item.ClosedAt = now
		}
		if !m.closedHas(id) {
			m.Closed = append(m.Closed, item)
		}
	}
	m.Active = kept
	if m.connDetailID == id {
		m.connDetailID = ""
	}
}

func (m *Model) closedHas(id string) bool {
	for _, item := range m.Closed {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (m *Model) markAllClosed(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	for _, item := range m.Active {
		item.Closed = true
		if item.ClosedAt.IsZero() {
			item.ClosedAt = now
		}
		if !m.closedHas(item.ID) {
			m.Closed = append(m.Closed, item)
		}
	}
	m.Active = nil
}

func (m *Model) clearClosed() {
	m.Closed = nil
	if m.connTab == 1 {
		m.connDetailID = ""
		m.connSelected = -1
	}
}

// ApplyConnections replaces the active set with a core snapshot. Flows
// that disappeared are appended to the closed list. Speeds are the byte
// delta since the previous snapshot, divided by the elapsed seconds.
func (m *Model) ApplyConnections(active []Connection, up, down int64, now time.Time) {
	if m == nil {
		return
	}
	if now.IsZero() {
		now = time.Now()
	}
	prev := make(map[string]Connection, len(m.Active))
	for _, item := range m.Active {
		prev[item.ID] = item
	}
	var dt float64
	if !m.connSynced.IsZero() {
		dt = now.Sub(m.connSynced).Seconds()
	}
	seen := make(map[string]struct{}, len(active))
	for i := range active {
		seen[active[i].ID] = struct{}{}
		active[i].Closed = false
		if old, ok := prev[active[i].ID]; ok && dt > 0 {
			active[i].UploadSpeed = perSec(active[i].Upload-old.Upload, dt)
			active[i].DownloadSpeed = perSec(active[i].Download-old.Download, dt)
		} else if old, ok := prev[active[i].ID]; ok {
			active[i].UploadSpeed = old.UploadSpeed
			active[i].DownloadSpeed = old.DownloadSpeed
		}
	}
	for _, item := range m.Active {
		if _, ok := seen[item.ID]; ok {
			continue
		}
		item.Closed = true
		if item.ClosedAt.IsZero() {
			item.ClosedAt = now
		}
		if !m.closedHas(item.ID) {
			m.Closed = append(m.Closed, item)
		}
	}
	if len(m.Closed) > 0 && len(seen) > 0 {
		kept := m.Closed[:0]
		for _, item := range m.Closed {
			if _, still := seen[item.ID]; still {
				continue
			}
			kept = append(kept, item)
		}
		m.Closed = kept
	}
	m.Active = active
	m.ConnUp = up
	m.ConnDown = down
	m.connSynced = now
}

func perSec(delta int64, dt float64) int64 {
	if delta <= 0 || dt <= 0 {
		return 0
	}
	return int64(float64(delta) / dt)
}

func filterConnections(items []Connection, query string, opt SearchOpts) []Connection {
	query = strings.TrimSpace(query)
	if query == "" {
		return append([]Connection(nil), items...)
	}
	re, ok := connRegexp(query, opt)
	if !ok {
		return nil
	}
	var out []Connection
	for _, item := range items {
		if connMatch(item, query, opt, re) {
			out = append(out, item)
		}
	}
	return out
}

func connRegexp(query string, opt SearchOpts) (*regexp.Regexp, bool) {
	if !opt.Regex {
		return nil, true
	}
	pat := query
	if !opt.MatchCase {
		pat = "(?i)" + pat
	}
	re, err := regexp.Compile(pat)
	if err != nil {
		return nil, false
	}
	return re, true
}

func connRegexInvalid(query string, opt SearchOpts) bool {
	query = strings.TrimSpace(query)
	if query == "" || !opt.Regex {
		return false
	}
	_, ok := connRegexp(query, opt)
	return !ok
}

func connMatch(item Connection, query string, opt SearchOpts, re *regexp.Regexp) bool {
	for _, field := range []string{item.Host, item.DestinationIP, item.Process} {
		if field != "" && nameMatch(field, query, opt, re) {
			return true
		}
	}
	return false
}

func sortConnections(items []Connection, order string, table bool, by mygoui.SortOrder) []Connection {
	out := append([]Connection(nil), items...)
	if table {
		if by.Column == "" {
			return out
		}
		slices.SortStableFunc(out, func(a, b Connection) int {
			cmp := compareConnection(a, b, by.Column)
			if by.Descending {
				return -cmp
			}
			return cmp
		})
		return out
	}
	slices.SortStableFunc(out, func(a, b Connection) int {
		switch order {
		case "uploadSpeed":
			return cmpInt64(b.UploadSpeed, a.UploadSpeed)
		case "downloadSpeed":
			return cmpInt64(b.DownloadSpeed, a.DownloadSpeed)
		default:
			return b.Start.Compare(a.Start)
		}
	})
	return out
}

func compareConnection(a, b Connection, column string) int {
	switch column {
	case "download":
		return cmpInt64(a.Download, b.Download)
	case "upload":
		return cmpInt64(a.Upload, b.Upload)
	case "dlSpeed":
		return cmpInt64(a.DownloadSpeed, b.DownloadSpeed)
	case "ulSpeed":
		return cmpInt64(a.UploadSpeed, b.UploadSpeed)
	case "time":
		return a.Start.Compare(b.Start)
	case "host":
		return strings.Compare(a.Endpoint(), b.Endpoint())
	case "chains":
		return strings.Compare(a.ChainLabel(), b.ChainLabel())
	case "rule":
		return strings.Compare(a.RuleLabel(), b.RuleLabel())
	case "process":
		return strings.Compare(a.Process, b.Process)
	case "source":
		return strings.Compare(a.SourceLabel(), b.SourceLabel())
	case "destination":
		return strings.Compare(a.DestinationLabel(), b.DestinationLabel())
	case "type":
		return strings.Compare(a.TypeLabel(), b.TypeLabel())
	default:
		return 0
	}
}

func cmpInt64(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func connOrderLabel(id string) string {
	switch id {
	case "uploadSpeed":
		return "上传速度"
	case "downloadSpeed":
		return "下载速度"
	default:
		return "默认"
	}
}

func connOrderID(label string) string {
	switch label {
	case "上传速度":
		return "uploadSpeed"
	case "下载速度":
		return "downloadSpeed"
	default:
		return "default"
	}
}

// Endpoint is the host the row shows: host, else destination IP, else the
// remote address, with the destination port.
func (c Connection) Endpoint() string {
	host := c.Host
	if host == "" {
		host = c.DestinationIP
	}
	if host == "" {
		host = c.RemoteDest
	}
	if host == "" {
		return "-"
	}
	if c.DestPort > 0 {
		return net.JoinHostPort(host, strconv.Itoa(int(c.DestPort)))
	}
	return host
}

// RuleLabel is rule or rule(payload), matching the Verge row.
func (c Connection) RuleLabel() string {
	if c.Rule == "" {
		return ""
	}
	if c.RulePayload == "" {
		return c.Rule
	}
	return c.Rule + "(" + c.RulePayload + ")"
}

// ChainLabel joins the chain from the outside in, as Verge reverses it.
func (c Connection) ChainLabel() string {
	if len(c.Chains) == 0 {
		return ""
	}
	parts := append([]string(nil), c.Chains...)
	slices.Reverse(parts)
	return strings.Join(parts, " / ")
}

// Tags are the list-row chips: network, type, process, chains, rule.
func (c Connection) Tags() []string {
	var out []string
	if c.Network != "" {
		out = append(out, c.Network)
	}
	if c.Type != "" {
		out = append(out, c.Type)
	}
	if c.Process != "" {
		out = append(out, c.Process)
	}
	if label := c.ChainLabel(); label != "" {
		out = append(out, label)
	}
	if label := c.RuleLabel(); label != "" {
		out = append(out, label)
	}
	return out
}

func (c Connection) SourceLabel() string {
	if c.SourceIP == "" && c.SourcePort == 0 {
		return ""
	}
	return net.JoinHostPort(c.SourceIP, strconv.Itoa(int(c.SourcePort)))
}

func (c Connection) DestinationLabel() string {
	if c.DestinationIP != "" {
		return c.DestinationIP
	}
	return c.RemoteDest
}

func (c Connection) TypeLabel() string {
	switch {
	case c.Type != "" && c.Network != "":
		return c.Type + "(" + c.Network + ")"
	case c.Type != "":
		return c.Type
	default:
		return c.Network
	}
}

func (c Connection) columnText(id string, now time.Time) string {
	switch id {
	case "host":
		return c.Endpoint()
	case "download":
		return FormatTraffic(c.Download)
	case "upload":
		return FormatTraffic(c.Upload)
	case "dlSpeed":
		return FormatSpeed(c.DownloadSpeed)
	case "ulSpeed":
		return FormatSpeed(c.UploadSpeed)
	case "chains":
		text := c.ChainLabel()
		if text == "" {
			return "-"
		}
		return text
	case "rule":
		text := c.RuleLabel()
		if text == "" {
			return "-"
		}
		return text
	case "process":
		if c.Process == "" {
			return "-"
		}
		return c.Process
	case "time":
		return RelativeSince(now, c.Start)
	case "source":
		text := c.SourceLabel()
		if text == "" {
			return "-"
		}
		return text
	case "destination":
		text := c.DestinationLabel()
		if text == "" {
			return "-"
		}
		return text
	case "type":
		text := c.TypeLabel()
		if text == "" {
			return "-"
		}
		return text
	default:
		return ""
	}
}

type connField struct {
	label string
	value string
}

func (c Connection) detailFields(now time.Time) []connField {
	process := c.Process
	if c.ProcessPath != "" {
		if process != "" {
			process += "(" + c.ProcessPath + ")"
		} else {
			process = c.ProcessPath
		}
	}
	return []connField{
		{"主机", c.Endpoint()},
		{"下载量", FormatTraffic(c.Download)},
		{"上传量", FormatTraffic(c.Upload)},
		{"下载速度", FormatSpeed(c.DownloadSpeed)},
		{"上传速度", FormatSpeed(c.UploadSpeed)},
		{"链路", c.ChainLabel()},
		{"规则", c.RuleLabel()},
		{"进程", process},
		{"连接时间", RelativeSince(now, c.Start)},
		{"源地址", c.SourceLabel()},
		{"目标地址", c.DestinationLabel()},
		{"目标端口", strconv.Itoa(int(c.DestPort))},
		{"类型", c.TypeLabel()},
	}
}

// RelativeSince is the short Chinese age of start at now.
func RelativeSince(now, start time.Time) string {
	if start.IsZero() {
		return "-"
	}
	d := now.Sub(start)
	if d < 0 {
		d = 0
	}
	switch {
	case d < 10*time.Second:
		return "刚刚"
	case d < time.Minute:
		return strconv.Itoa(int(d.Seconds())) + " 秒前"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + " 分钟前"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d.Hours())) + " 小时前"
	default:
		return strconv.Itoa(int(d.Hours()/24)) + " 天前"
	}
}
