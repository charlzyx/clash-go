package ui

import mygoui "github.com/egoist/mygo/ui"

// Page is one entry in the sidebar. ID is the screenshot file prefix.
type Page struct {
	ID    string
	Path  string
	Title string
	Body  string
}

// Pages is the sidebar, in Clash Verge Rev order.
func Pages() []Page {
	return []Page{
		{ID: "home", Path: "/home", Title: "首页"},
		{ID: "proxies", Path: "/proxies", Title: "代理", Body: "节点和策略组会显示在这里。"},
		{ID: "subscriptions", Path: "/subscriptions", Title: "订阅", Body: "订阅配置会显示在这里。"},
		{ID: "connections", Path: "/connections", Title: "连接", Body: "活动连接会显示在这里。"},
		{ID: "rules", Path: "/rules", Title: "规则", Body: "分流规则会显示在这里。"},
		{ID: "logs", Path: "/logs", Title: "日志", Body: "内核日志会显示在这里。"},
		{ID: "test", Path: "/test", Title: "测试", Body: "连通性测试会显示在这里。"},
		{ID: "settings", Path: "/settings", Title: "设置", Body: "应用设置会显示在这里。"},
	}
}

// HomeState is everything the home page renders. Callers pass a snapshot;
// the page does not read the kernel.
type HomeState struct {
	Version     string
	Running     bool
	MixedPort   int
	Mode        string
	ConfigDir   string
	Error       string
	SystemProxy bool
	TUN         bool
}

// Model is the main window: a sidebar and the page it shows.
type Model struct {
	Router *mygoui.Router
	Home   HomeState
}

// NewModel returns a window opened on the home page.
func NewModel(home HomeState) *Model {
	return &Model{Router: mygoui.NewRouter("/home"), Home: home}
}

// View builds one frame of the main window.
func (m *Model) View(c *mygoui.Context) {
	if m == nil {
		return
	}
	t := c.Theme()
	c.Root().Background(t.Background)
	mygoui.Row(c).Fill().AlignItems(mygoui.Stretch).Children(func() {
		m.sidebar(c)
		mygoui.Column(c).Grow(1).Padding(28).Gap(8).Children(func() {
			m.Router.View(c, func(r *mygoui.Route) {
				path := m.Router.Path()
				if path == "/home" {
					r.Title("首页")
					HomePage(c, m.Home)
					return
				}
				page, ok := pageByPath(path)
				if !ok {
					r.Title("未找到")
					PlaceholderPage(c, "未找到", "这个页面还没有。")
					return
				}
				r.Title(page.Title)
				PlaceholderPage(c, page.Title, page.Body)
			})
		})
	})
}

func (m *Model) sidebar(c *mygoui.Context) {
	t := c.Theme()
	page := m.Router.Path()
	mygoui.Column(c).Width(200).FillHeight().Background(t.Surface).Children(func() {
		mygoui.Column(c).Padding(16, 14, 8, 14).Gap(2).Children(func() {
			mygoui.Text(c, "Clash Go").Bold().FontSize(16)
			mygoui.Text(c, sidebarStatus(m.Home)).FontSize(11).TextColor(t.TextMuted)
		})
		if mygoui.Sidebar(c, &page, func() {
			for _, item := range Pages() {
				mygoui.SidebarItem(c, item.Path, nil, item.Title)
			}
		}).Grow(1).Changed() {
			m.Router.Push(page)
		}
	})
}

func sidebarStatus(s HomeState) string {
	if s.Running {
		return "运行中"
	}
	return "已停止"
}

func pageByPath(path string) (Page, bool) {
	for _, p := range Pages() {
		if p.Path == path {
			return p, true
		}
	}
	return Page{}, false
}

// HomePage renders the home dashboard from s.
func HomePage(c *mygoui.Context, s HomeState) {
	t := c.Theme()
	mygoui.Column(c).Gap(16).Children(func() {
		mygoui.Text(c, "首页").FontSize(28).Bold()
		mygoui.Text(c, "内核").FontSize(13).TextColor(t.TextMuted)
		mygoui.Row(c).Gap(12).AlignItems(mygoui.Stretch).Children(func() {
			statCard(c, "内核版本", display(s.Version))
			statCard(c, "运行状态", runLabel(s.Running))
			statCard(c, "混合端口", portLabel(s.MixedPort))
		})
		if s.Error != "" {
			mygoui.Text(c, s.Error).TextColor(t.Danger)
		}
		mygoui.Column(c).Gap(4).Children(func() {
			mygoui.Text(c, "模式").FontSize(12).TextColor(t.TextMuted)
			mygoui.Text(c, display(s.Mode))
			mygoui.Text(c, "配置目录").FontSize(12).TextColor(t.TextMuted)
			mygoui.Text(c, display(s.ConfigDir))
			mygoui.Text(c, "系统代理 "+onOff(s.SystemProxy)).TextColor(t.TextMuted)
			mygoui.Text(c, "TUN "+onOff(s.TUN)).TextColor(t.TextMuted)
		})
	})
}

// PlaceholderPage renders a page that later milestones fill in.
func PlaceholderPage(c *mygoui.Context, title, body string) {
	t := c.Theme()
	mygoui.Column(c).Gap(8).Children(func() {
		mygoui.Text(c, title).FontSize(28).Bold()
		mygoui.Text(c, body).TextColor(t.TextMuted)
	})
}

func statCard(c *mygoui.Context, label, value string) {
	t := c.Theme()
	mygoui.Column(c).Grow(1).Padding(16).Gap(6).Radius(12).Background(t.Surface).Border(1, t.Border).Children(func() {
		mygoui.Text(c, label).FontSize(12).TextColor(t.TextMuted)
		mygoui.Text(c, value).FontSize(22).Bold()
	})
}

func runLabel(running bool) string {
	if running {
		return "运行中"
	}
	return "已停止"
}

func portLabel(port int) string {
	if port <= 0 {
		return "—"
	}
	return itoa(port)
}

func onOff(v bool) string {
	if v {
		return "开启"
	}
	return "关闭"
}

func display(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	neg := n < 0
	if neg {
		n = -n
	}
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
