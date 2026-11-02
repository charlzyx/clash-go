package ui

import (
	"errors"
	"testing"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

func connPage() (*Model, *mygoui.Tester) {
	m := NewModel(demoHome())
	m.Router.Replace("/connections")
	return m, mygoui.NewTester(m.View, WindowWidth, WindowHeight)
}

func sampleConn(id, host string) Connection {
	return Connection{
		ID: id, Host: host, DestPort: 443, Process: "Safari",
		Rule: "DOMAIN-SUFFIX", RulePayload: host, Network: "tcp", Type: "HTTPS",
		Download: 2048, Upload: 512, Start: time.Now().Add(-time.Minute),
	}
}

func TestConnectionsEmpty(t *testing.T) {
	_, tt := connPage()
	for _, text := range []string{"连接", "空空如也", "活跃", "已关闭", "关闭全部", "下载量", "上传量", "默认", "conn-search"} {
		if !tt.HasText(text) {
			t.Errorf("empty connections page missing %q\n%v", text, tt.Texts())
		}
	}
	if tt.HasText("活动连接会显示在这里。") {
		t.Fatal("connections is still the placeholder")
	}
}

func TestConnectionsActiveRowAndTabs(t *testing.T) {
	m, _ := connPage()
	m.Active = []Connection{sampleConn("a", "example.com")}
	m.Closed = []Connection{sampleConn("b", "old.example")}
	m.Closed[0].Closed = true
	m.Closed[0].Process = "curl"
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	for _, text := range []string{"example.com:443", "Safari", "DOMAIN-SUFFIX(example.com)", "下载 2.00 KB", "上传 512 B"} {
		if !tt.HasText(text) {
			t.Errorf("active row missing %q\n%v", text, tt.Texts())
		}
	}
	if tt.HasText("old.example") {
		t.Fatal("closed row is visible on the active tab")
	}
	if err := tt.Click("已关闭"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("old.example:443") || !tt.HasText("curl") {
		t.Fatalf("closed tab texts %q", tt.Texts())
	}
	if tt.HasText("example.com:443") {
		t.Fatal("active row stayed visible on the closed tab")
	}
	if err := tt.Click("old.example:443"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("连接详情") || !tt.HasText("已关闭") {
		t.Fatal("closed detail did not open")
	}
	if tt.HasText("关闭连接") {
		t.Fatal("closed detail offered 关闭连接")
	}
	if err := tt.Click("conn-detail-close"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("活跃"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("example.com:443"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("关闭连接") || !tt.HasText("活跃") {
		t.Fatalf("active detail texts %q", tt.Texts())
	}
}

func TestConnectionsCloseAll(t *testing.T) {
	m, _ := connPage()
	m.Active = []Connection{sampleConn("a", "example.com")}
	var calls int
	fail := true
	m.Hooks.CloseAllConnections = func() error {
		calls++
		if fail {
			return errors.New("boom")
		}
		return nil
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("关闭全部"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("CloseAll calls %d, want 1", calls)
	}
	if m.Notice != "" {
		t.Fatalf("close-all failure set notice %q", m.Notice)
	}
	if !tt.HasText("example.com:443") {
		t.Fatal("failed close-all removed the row")
	}
	fail = false
	if err := tt.Click("关闭全部"); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("CloseAll calls %d, want 2", calls)
	}
	if m.Notice != "" {
		t.Fatalf("close-all success set notice %q", m.Notice)
	}
	if tt.HasText("example.com:443") {
		t.Fatal("active row remained after close-all")
	}
	if err := tt.Click("已关闭"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("example.com:443") {
		t.Fatal("closed tab did not keep the flow")
	}
}

func TestConnectionsSearchFilters(t *testing.T) {
	m, _ := connPage()
	m.Active = []Connection{sampleConn("a", "example.com"), sampleConn("b", "github.com")}
	m.Active[1].Process = "curl"
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("example.com:443") || !tt.HasText("github.com:443") {
		t.Fatalf("both rows should show: %q", tt.Texts())
	}
	if err := tt.Click("conn-search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("example")
	if !tt.HasText("example.com:443") || tt.HasText("github.com:443") {
		t.Fatalf("search texts %q", tt.Texts())
	}
	m.ConnQuery = ""
	tt = mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("github.com:443") {
		t.Fatal("clearing the query did not restore the row")
	}
}

func TestConnectionsLayoutHidesSort(t *testing.T) {
	m, _ := connPage()
	m.Active = []Connection{sampleConn("a", "example.com")}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("默认") {
		t.Fatal("list layout should show the sort menu")
	}
	if err := tt.Click("表格"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("默认") {
		t.Fatal("table layout still shows the list sort menu")
	}
	if !tt.HasText("列") || m.ConnLayout != "table" {
		t.Fatalf("layout %q texts %q", m.ConnLayout, tt.Texts())
	}
	if err := tt.Click("列表"); err != nil {
		t.Fatal(err)
	}
	if m.ConnLayout != "list" || !tt.HasText("默认") {
		t.Fatal("list layout did not come back")
	}
}

func TestConnectionsScrollShowsTopButton(t *testing.T) {
	m, _ := connPage()
	base := time.Unix(2_000_000, 0)
	for i := 0; i < 40; i++ {
		item := sampleConn(itoa(i), "host"+itoa(i)+".example")
		item.Start = base.Add(time.Duration(40-i) * time.Second)
		m.Active = append(m.Active, item)
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if tt.HasText("回到顶部") {
		t.Fatal("scroll-to-top is visible before scrolling")
	}
	row, ok := tt.Find("host0.example:443")
	if !ok {
		t.Fatalf("missing first row in %v", tt.Texts())
	}
	tt.Scroll(row.X+20, row.Y+row.H/2, 0, 240)
	if m.connScroll.Y <= connScrollTop || !tt.HasText("回到顶部") {
		t.Fatalf("after scroll Y=%v texts %v", m.connScroll.Y, tt.Texts())
	}
	if err := tt.Click("回到顶部"); err != nil {
		t.Fatal(err)
	}
	if m.connScroll.Y > 1 {
		t.Fatalf("scroll Y %v after jumping to the top", m.connScroll.Y)
	}
}

func TestConnectionFilterAndSort(t *testing.T) {
	items := []Connection{
		{ID: "1", Host: "Example.COM", Process: "Safari", DestinationIP: "203.0.113.5", Start: time.Unix(10, 0), UploadSpeed: 10, DownloadSpeed: 1},
		{ID: "2", Host: "github.com", Process: "curl", Start: time.Unix(30, 0), UploadSpeed: 3, DownloadSpeed: 90},
		{ID: "3", Host: "example.org", Process: "ExampleHelper", Start: time.Unix(20, 0), UploadSpeed: 50, DownloadSpeed: 20},
	}
	if got := filterConnections(items, "  ", SearchOpts{}); len(got) != 3 {
		t.Fatalf("blank query filtered to %d", len(got))
	}
	got := filterConnections(items, "example", SearchOpts{})
	if len(got) != 2 || got[0].ID != "1" || got[1].ID != "3" {
		t.Fatalf("case-insensitive host filter %#v", ids(got))
	}
	got = filterConnections(items, "example", SearchOpts{MatchCase: true})
	if len(got) != 1 || got[0].ID != "3" {
		t.Fatalf("match-case %#v", ids(got))
	}
	got = filterConnections(items, "Safari", SearchOpts{WholeWord: true})
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("whole word %#v", ids(got))
	}
	if got = filterConnections(items, "Saf", SearchOpts{WholeWord: true}); len(got) != 0 {
		t.Fatalf("whole word prefix %#v", ids(got))
	}
	got = filterConnections(items, `example\.com`, SearchOpts{Regex: true})
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("regex %#v", ids(got))
	}
	if got = filterConnections(items, "[", SearchOpts{Regex: true}); got != nil {
		t.Fatalf("bad regex returned %#v", ids(got))
	}
	ordered := sortConnections(items, "default", false, mygoui.SortOrder{})
	if ids(ordered)[0] != "2" || ids(ordered)[2] != "1" {
		t.Fatalf("default order %#v", ids(ordered))
	}
	ordered = sortConnections(items, "uploadSpeed", false, mygoui.SortOrder{})
	if ids(ordered)[0] != "3" {
		t.Fatalf("upload speed order %#v", ids(ordered))
	}
	ordered = sortConnections(items, "downloadSpeed", false, mygoui.SortOrder{})
	if ids(ordered)[0] != "2" {
		t.Fatalf("download speed order %#v", ids(ordered))
	}
}

func TestApplyConnectionsMovesClosedAndSpeeds(t *testing.T) {
	m := NewModel(demoHome())
	t0 := time.Unix(1000, 0)
	m.ApplyConnections([]Connection{{ID: "a", Host: "example.com", Upload: 100, Download: 1000}}, 100, 1000, t0)
	t1 := t0.Add(2 * time.Second)
	m.ApplyConnections([]Connection{{ID: "a", Host: "example.com", Upload: 300, Download: 1000}}, 300, 1000, t1)
	if m.Active[0].UploadSpeed != 100 || m.Active[0].DownloadSpeed != 0 {
		t.Fatalf("speeds up %d down %d", m.Active[0].UploadSpeed, m.Active[0].DownloadSpeed)
	}
	m.ApplyConnections(nil, 300, 1000, t1.Add(time.Second))
	if len(m.Active) != 0 || len(m.Closed) != 1 || !m.Closed[0].Closed || m.Closed[0].ID != "a" {
		t.Fatalf("active %#v closed %#v", m.Active, m.Closed)
	}
	m.ApplyConnections([]Connection{{ID: "a", Host: "example.com"}}, 300, 1000, t1.Add(2*time.Second))
	if len(m.Active) != 1 || len(m.Closed) != 0 {
		t.Fatalf("revived flow stayed closed: active %d closed %d", len(m.Active), len(m.Closed))
	}
}

func TestRelativeSince(t *testing.T) {
	now := time.Unix(10_000, 0)
	if RelativeSince(now, time.Time{}) != "-" {
		t.Fatal("zero")
	}
	if RelativeSince(now, now.Add(-3*time.Second)) != "刚刚" {
		t.Fatal("just now")
	}
	if RelativeSince(now, now.Add(-30*time.Second)) != "30 秒前" {
		t.Fatal("seconds")
	}
	if RelativeSince(now, now.Add(-5*time.Minute)) != "5 分钟前" {
		t.Fatal("minutes")
	}
}

func ids(items []Connection) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.ID
	}
	return out
}
