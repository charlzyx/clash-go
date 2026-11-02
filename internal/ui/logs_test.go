package ui

import (
	"testing"

	mygoui "github.com/egoist/mygo/ui"
)

func logPage() (*Model, *mygoui.Tester) {
	m := NewModel(demoHome())
	m.Router.Replace("/logs")
	return m, mygoui.NewTester(m.View, WindowWidth, WindowHeight)
}

func TestLogsEmpty(t *testing.T) {
	_, tt := logPage()
	for _, text := range []string{"日志", "空空如也", "暂停", "顺序", "清除", "ALL", "log-search"} {
		if !tt.HasText(text) {
			t.Errorf("empty logs page missing %q\n%v", text, tt.Texts())
		}
	}
	if tt.HasText("内核日志会显示在这里。") {
		t.Fatal("logs is still the placeholder")
	}
	bar, ok := tt.Find("log-toolbar")
	if !ok {
		t.Fatal("missing log toolbar")
	}
	if bar.H < 38 || bar.H > 42 {
		t.Fatalf("toolbar height %v, want about 39", bar.H)
	}
}

func TestLogsRowsVisible(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{
		{ID: 1, Time: "11-02 09:30:01", Type: "info", Payload: "hello-from-core"},
		{ID: 2, Time: "11-02 09:30:02", Type: "warning", Payload: "careful-now"},
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	for _, text := range []string{"11-02 09:30:01", "INFO", "hello-from-core", "WARNING", "careful-now"} {
		if !tt.HasText(text) {
			t.Errorf("log row missing %q\n%v", text, tt.Texts())
		}
	}
	if tt.HasText("空空如也") {
		t.Fatal("empty state stayed up with rows")
	}
}

func TestLogsLevelFilter(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{
		{ID: 1, Time: "11-02 09:30:01", Type: "info", Payload: "info-payload"},
		{ID: 2, Time: "11-02 09:30:02", Type: "warning", Payload: "warn-payload"},
		{ID: 3, Time: "11-02 09:30:03", Type: "error", Payload: "error-payload"},
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("log-filter"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("WARN"); err != nil {
		t.Fatal(err)
	}
	if m.ClashLog.LogFilter != "warn" {
		t.Fatalf("filter %q, want warn", m.ClashLog.LogFilter)
	}
	if !tt.HasText("warn-payload") {
		t.Fatalf("warning row hidden: %q", tt.Texts())
	}
	if tt.HasText("info-payload") || tt.HasText("error-payload") {
		t.Fatalf("level filter kept other rows: %q", tt.Texts())
	}
}

func TestLogsSearchFilters(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{
		{ID: 1, Time: "11-02 08:15:00", Type: "info", Payload: "alpha-line"},
		{ID: 2, Time: "11-02 09:15:00", Type: "error", Payload: "beta-line"},
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("log-search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("alpha")
	if !tt.HasText("alpha-line") || tt.HasText("beta-line") {
		t.Fatalf("payload search texts %q", tt.Texts())
	}

	m.LogQuery = ""
	tt = mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("log-search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("08:15")
	if !tt.HasText("alpha-line") || tt.HasText("beta-line") {
		t.Fatalf("time search texts %q", tt.Texts())
	}

	m.LogQuery = ""
	tt = mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("log-search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("ERROR")
	if !tt.HasText("beta-line") || tt.HasText("alpha-line") {
		t.Fatalf("type search texts %q", tt.Texts())
	}
}

func TestLogsPauseResume(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{{ID: 1, Time: "11-02 09:30:01", Type: "info", Payload: "kept-line"}}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !m.ClashLog.Enable || !tt.HasText("暂停") {
		t.Fatalf("enable %v texts %q", m.ClashLog.Enable, tt.Texts())
	}
	if err := tt.Click("log-pause"); err != nil {
		t.Fatal(err)
	}
	if m.ClashLog.Enable {
		t.Fatal("pause did not clear clashLog.enable")
	}
	if !tt.HasText("继续") || tt.HasText("暂停") {
		t.Fatalf("resume title missing: %q", tt.Texts())
	}
	if !tt.HasText("kept-line") {
		t.Fatal("pause cleared the rows")
	}
	m.AppendLogs([]LogItem{{Type: "info", Payload: "dropped-while-paused"}})
	if len(m.Logs) != 1 || m.Logs[0].Payload != "kept-line" {
		t.Fatalf("paused collection recorded %+v", m.Logs)
	}
	if err := tt.Click("log-pause"); err != nil {
		t.Fatal(err)
	}
	if !m.ClashLog.Enable || !tt.HasText("暂停") {
		t.Fatalf("resume enable %v texts %q", m.ClashLog.Enable, tt.Texts())
	}
	m.AppendLogs([]LogItem{{Type: "info", Payload: "after-resume"}})
	tt = mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("after-resume") {
		t.Fatal("resume did not record a line")
	}
}

func TestLogsClearInvokesHook(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{{ID: 1, Time: "11-02 09:30:01", Type: "info", Payload: "wipe-me"}}
	var calls int
	m.Hooks.ClearLogs = func() { calls++ }
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if err := tt.Click("清除"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("clear calls %d, want 1", calls)
	}
	if len(m.Logs) != 0 {
		t.Fatalf("rows left after clear: %+v", m.Logs)
	}
	if tt.HasText("wipe-me") || !tt.HasText("空空如也") {
		t.Fatalf("after clear texts %q", tt.Texts())
	}
}

func TestLogsOrderReverses(t *testing.T) {
	m, _ := logPage()
	m.Logs = []LogItem{
		{ID: 1, Time: "11-02 09:30:01", Type: "info", Payload: "older-line"},
		{ID: 2, Time: "11-02 09:30:02", Type: "info", Payload: "newer-line"},
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	older, ok1 := tt.Find("older-line")
	newer, ok2 := tt.Find("newer-line")
	if !ok1 || !ok2 {
		t.Fatalf("missing rows %v", tt.Texts())
	}
	if older.Y >= newer.Y {
		t.Fatalf("asc places newer above older: older %v newer %v", older, newer)
	}
	if err := tt.Click("log-order"); err != nil {
		t.Fatal(err)
	}
	if m.ClashLog.LogOrder != "desc" || !tt.HasText("倒序") {
		t.Fatalf("order %q texts %q", m.ClashLog.LogOrder, tt.Texts())
	}
	older, ok1 = tt.Find("older-line")
	newer, ok2 = tt.Find("newer-line")
	if !ok1 || !ok2 || newer.Y >= older.Y {
		t.Fatalf("desc did not reverse: older %v newer %v ok %v %v", older, newer, ok1, ok2)
	}
}

func TestLogsFollowsNewest(t *testing.T) {
	m, _ := logPage()
	for i := 0; i < 80; i++ {
		m.Logs = append(m.Logs, LogItem{
			ID: uint64(i + 1), Time: "11-02 10:00:00", Type: "info", Payload: "line-" + itoa(i),
		})
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("line-79") {
		t.Fatal("ascending log did not keep the newest line in view")
	}
	if tt.HasText("line-0") {
		t.Fatal("oldest line stayed in view")
	}
}

func TestFilterLogs(t *testing.T) {
	items := []LogItem{
		{Time: "11-02 08:15:00", Type: "info", Payload: "Alpha"},
		{Time: "11-02 09:15:00", Type: "warning", Payload: "beta"},
		{Time: "11-02 10:15:00", Type: "error", Payload: "gamma"},
		{Time: "11-02 11:15:00", Type: "debug", Payload: "delta"},
	}
	if got := filterLogs(items, "all", "", SearchOpts{}); len(got) != 4 {
		t.Fatalf("all %d", len(got))
	}
	got := filterLogs(items, "warn", "", SearchOpts{})
	if len(got) != 1 || got[0].Payload != "beta" {
		t.Fatalf("warn %+v", got)
	}
	got = filterLogs(items, "err", "", SearchOpts{})
	if len(got) != 1 || got[0].Type != "error" {
		t.Fatalf("err %+v", got)
	}
	got = filterLogs(items, "info", "", SearchOpts{})
	if len(got) != 1 || got[0].Payload != "Alpha" {
		t.Fatalf("info %+v", got)
	}
	got = filterLogs(items, "all", "ALPHA", SearchOpts{})
	if len(got) != 1 || got[0].Payload != "Alpha" {
		t.Fatalf("case-insensitive search %+v", got)
	}
	got = filterLogs(items, "all", "ALPHA", SearchOpts{MatchCase: true})
	if len(got) != 0 {
		t.Fatalf("match-case against lowercased text %+v", got)
	}
	got = filterLogs(items, "all", "08:15", SearchOpts{})
	if len(got) != 1 || got[0].Payload != "Alpha" {
		t.Fatalf("time search %+v", got)
	}
	got = filterLogs(items, "all", "warning", SearchOpts{})
	if len(got) != 1 || got[0].Payload != "beta" {
		t.Fatalf("type search %+v", got)
	}
	if got = filterLogs(items, "all", "[", SearchOpts{Regex: true}); got != nil {
		t.Fatalf("invalid regex %+v", got)
	}
	rev := reverseLogs(filterLogs(items, "all", "", SearchOpts{}))
	if rev[0].Payload != "delta" || rev[3].Payload != "Alpha" {
		t.Fatalf("reverse %+v", rev)
	}
}

func TestAppendLogsPauseAndCap(t *testing.T) {
	m := NewModel(demoHome())
	m.ClashLog.Enable = false
	m.AppendLogs([]LogItem{{Type: "info", Payload: "nope"}})
	if len(m.Logs) != 0 {
		t.Fatalf("paused append kept %+v", m.Logs)
	}
	m.ClashLog.Enable = true
	batch := make([]LogItem, maxLogRows+5)
	for i := range batch {
		batch[i] = LogItem{Type: "info", Payload: itoa(i)}
	}
	m.AppendLogs(batch)
	if len(m.Logs) != maxLogRows {
		t.Fatalf("len %d, want %d", len(m.Logs), maxLogRows)
	}
	if m.Logs[0].Payload != "5" || m.Logs[len(m.Logs)-1].Payload != itoa(maxLogRows+4) {
		t.Fatalf("cap kept %q .. %q", m.Logs[0].Payload, m.Logs[len(m.Logs)-1].Payload)
	}
	if m.Logs[0].ID == 0 || m.Logs[0].ID == m.Logs[1].ID {
		t.Fatalf("ids %#v %#v", m.Logs[0].ID, m.Logs[1].ID)
	}
}
