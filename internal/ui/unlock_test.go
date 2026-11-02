package ui

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestDedupeUnlockPrefersResolvedOverPending(t *testing.T) {
	got := DedupeUnlockItems([]UnlockItem{
		{Name: "Netflix", Status: "Pending", Region: "US"},
		{Name: " netflix ", Status: "Yes", CheckTime: "2026-11-02 01:02:03"},
		{Name: "Disney+", Status: "No", Region: "JP"},
		{Name: "disney+", Status: "Pending", CheckTime: "2026-11-02 04:05:06"},
	})
	if len(got) != 2 {
		t.Fatalf("len %d %+v", len(got), got)
	}
	if got[0].Name != " netflix " || got[0].Status != "Yes" || got[0].Region != "US" || got[0].CheckTime != "2026-11-02 01:02:03" {
		t.Fatalf("netflix %+v", got[0])
	}
	if got[1].Status != "No" || got[1].Region != "JP" || got[1].CheckTime != "2026-11-02 04:05:06" {
		t.Fatalf("disney %+v", got[1])
	}
}

func TestMergeUnlockDropsUnknownAndSorts(t *testing.T) {
	merged := MergeUnlockItems(DefaultUnlockItems(), []UnlockItem{
		{Name: "netflix", Status: "Yes", Region: "US", CheckTime: "2026-11-02 08:00:00"},
		{Name: "Not A Service", Status: "Yes"},
		{Name: "Netflix", Status: "Pending"},
	})
	if len(merged) != len(DefaultUnlockItems()) {
		t.Fatalf("len %d", len(merged))
	}
	var netflix UnlockItem
	for _, item := range merged {
		if item.Name == "Not A Service" {
			t.Fatal("unknown name was kept")
		}
		if strings.EqualFold(item.Name, "netflix") {
			netflix = item
		}
	}
	if netflix.Name != "Netflix" || netflix.Status != "Yes" || netflix.Region != "US" {
		t.Fatalf("netflix %+v", netflix)
	}
	sorted := SortUnlockItems(merged)
	for i := range merged {
		if merged[i].Name != sorted[i].Name {
			t.Fatalf("merge did not sort: %+v", namesOf(merged))
		}
	}
}

func TestUnlockStatusTable(t *testing.T) {
	cases := []struct {
		status, label, border, icon string
		tone                        Tone
		bold                        bool
	}{
		{"Pending", "待检测", "divider", "pending", ToneMuted, false},
		{"Yes", "支持", "success", "yes", ToneSuccess, true},
		{"No", "不支持", "error", "no", ToneError, true},
		{"Soon", "Soon", "warning", "soon", ToneWarning, true},
		{"Risky", "存疑", "warning", "help", ToneWarning, true},
		{"Failed", "测试失败", "error", "help", ToneError, true},
		{"Failed (Network Connection)", "测试失败（网络连接问题）", "error", "help", ToneError, true},
		{"Completed", "检测完成", "info", "help", ToneInfo, true},
		{"Disallowed ISP", "不允许的 ISP", "divider", "help", ToneError, true},
		{"Blocked", "Blocked", "divider", "help", ToneError, true},
		{"Unsupported Country/Region", "不支持的国家/地区", "divider", "help", ToneError, true},
		{"Originals Only", "仅限原创", "divider", "help", ToneMuted, true},
		{"No (IP Banned By Disney+)", "不支持（IP 被 Disney+ 禁止）", "divider", "help", ToneMuted, true},
	}
	for _, tc := range cases {
		got := UnlockStatusStyleOf(tc.status)
		if got.Label != tc.label || got.Border != tc.border || got.Icon != tc.icon || got.Tone != tc.tone || got.Bold != tc.bold {
			t.Errorf("%s -> %+v", tc.status, got)
		}
	}
}

func TestUnlockPageShowsCardsAndEmpty(t *testing.T) {
	m := NewModel(demoHome())
	m.Hooks.RefreshUnlock = func() {
		m.Unlock = []UnlockItem{
			{Name: "Netflix", Status: "Yes", Region: "US", CheckTime: "2026-11-02 08:00:00"},
			{Name: "Disney+", Status: "No", CheckTime: "2026-11-02 08:00:01"},
			{Name: "Muse AI", Status: "Risky", Region: "US", CheckTime: "2026-11-02 08:00:02"},
		}
	}
	m.Router.Replace("/unlock")
	tt := NewPageTester(m)
	for _, text := range []string{"解锁测试", "全部测试", "Netflix", "支持", "US", "Disney+", "不支持", "Muse AI", "存疑", "2026-11-02 08:00:00"} {
		if !tt.HasText(text) {
			t.Errorf("missing %q\n%s", text, join(tt.Texts()))
		}
	}
	if tt.HasText("连通性测试会显示在这里。") {
		t.Fatal("unlock is still the placeholder")
	}
	if _, ok := tt.Find("unlock-Netflix"); !ok {
		t.Fatal("card is not clickable")
	}
	if tt.HasText("空空如也") {
		t.Fatal("empty state shows with cards")
	}

	m.Hooks.RefreshUnlock = func() { m.Unlock = nil }
	m.Router.Replace("/home")
	tt = NewPageTester(m)
	if err := tt.Click("测 试"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("空空如也") || !tt.HasText("解锁测试") {
		t.Fatalf("empty\n%s", join(tt.Texts()))
	}
	if _, ok := tt.Find("unlock-empty"); !ok {
		t.Fatal("empty region missing")
	}
}

func TestUnlockColumnsFollowWindow(t *testing.T) {
	if UnlockColumns(940) != 3 || UnlockColumns(700) != 2 || UnlockColumns(520) != 1 {
		t.Fatalf("cols %d %d %d", UnlockColumns(940), UnlockColumns(700), UnlockColumns(520))
	}
	m := NewModel(demoHome())
	m.Hooks.RefreshUnlock = func() {
		m.Unlock = []UnlockItem{
			{Name: "A", Status: "Yes"},
			{Name: "B", Status: "No"},
			{Name: "C", Status: "Risky"},
			{Name: "D", Status: "Pending"},
		}
	}
	m.Router.Replace("/unlock")
	wide := renderModel(m, WindowWidth, WindowHeight, false)
	a, aOK := wide.Find("unlock-A")
	b, bOK := wide.Find("unlock-B")
	c, cOK := wide.Find("unlock-C")
	d, dOK := wide.Find("unlock-D")
	if !aOK || !bOK || !cOK || !dOK {
		t.Fatalf("cards %#v %#v %#v %#v", a, b, c, d)
	}
	if !sameRow(a.Y, b.Y) || !sameRow(a.Y, c.Y) || sameRow(a.Y, d.Y) {
		t.Fatalf("wide rows a=%v b=%v c=%v d=%v", a, b, c, d)
	}
	narrow := renderModel(m, 520, WindowHeight, false)
	na, _ := narrow.Find("unlock-A")
	nb, _ := narrow.Find("unlock-B")
	if sameRow(na.Y, nb.Y) {
		t.Fatalf("narrow still shares a row: %#v %#v", na, nb)
	}
}

func TestUnlockRetestAndFailureNotice(t *testing.T) {
	m := NewModel(demoHome())
	m.Hooks.RefreshUnlock = func() {
		m.Unlock = []UnlockItem{{Name: "Netflix", Status: "Pending"}}
	}
	var saved []UnlockItem
	m.Hooks.SavePrefs = func(p Prefs) error {
		saved = append([]UnlockItem(nil), p.UnlockResults...)
		return nil
	}
	m.Router.Replace("/unlock")
	tt := NewPageTester(m)
	m.Hooks.CheckUnlock = func(name string, onItem func(UnlockItem)) ([]UnlockItem, error) {
		if name != "Netflix" {
			t.Fatalf("name %q", name)
		}
		return nil, errors.New("boom")
	}
	if err := tt.Click("unlock-Netflix"); err != nil {
		t.Fatal(err)
	}
	if m.Notice != "Netflix 检测失败" || !tt.HasText("Netflix 检测失败") {
		t.Fatalf("notice %q\n%s", m.Notice, join(tt.Texts()))
	}
	if saved != nil {
		t.Fatalf("failure saved %+v", saved)
	}
	if m.Unlock[0].Status != "Pending" {
		t.Fatalf("status changed on failure %+v", m.Unlock[0])
	}

	m.Hooks.CheckUnlock = func(name string, onItem func(UnlockItem)) ([]UnlockItem, error) {
		item := UnlockItem{Name: name, Status: "Yes", Region: "US", CheckTime: "2026-11-02 09:00:00"}
		if onItem != nil {
			onItem(item)
		}
		return []UnlockItem{item}, nil
	}
	if err := tt.Click("unlock-Netflix"); err != nil {
		t.Fatal(err)
	}
	if m.Unlock[0].Status != "Yes" || m.Notice != "" {
		t.Fatalf("after retest %+v notice %q", m.Unlock[0], m.Notice)
	}
	if len(saved) != 1 || saved[0].Status != "Yes" || saved[0].Region != "US" {
		t.Fatalf("cache %+v", saved)
	}
	raw, err := json.Marshal(m.Prefs)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"clash_verge_unlock_results"`) {
		t.Fatalf("prefs json %s", raw)
	}
}

func TestUnlockTestAllProgressAndTimeout(t *testing.T) {
	m := NewModel(demoHome())
	m.Hooks.RefreshUnlock = func() {
		m.Unlock = []UnlockItem{
			{Name: "Netflix", Status: "Pending"},
			{Name: "Disney+", Status: "Pending"},
		}
	}
	m.Router.Replace("/unlock")
	tt := NewPageTester(m)
	m.unlockAll = true
	m.unlockLoading["Netflix"] = true
	tt = NewPageTester(m)
	if !tt.HasText("测试中... 1/2") {
		t.Fatalf("progress\n%s", join(tt.Texts()))
	}
	calls := 0
	if err := tt.Click("unlock-test-all"); err != nil {
		t.Fatal(err)
	}
	if calls != 0 {
		t.Fatal("testing button still started a run")
	}

	m.unlockAll = false
	m.unlockLoading = map[string]bool{}
	m.Hooks.CheckUnlock = func(name string, onItem func(UnlockItem)) ([]UnlockItem, error) {
		calls++
		if name != "" {
			t.Fatalf("all-test name %q", name)
		}
		return nil, errors.New("timeout")
	}
	tt = NewPageTester(m)
	if err := tt.Click("unlock-test-all"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || m.Notice != "检测超时或失败" || m.unlockAll {
		t.Fatalf("calls %d notice %q all %v", calls, m.Notice, m.unlockAll)
	}

	m.Hooks.CheckUnlock = func(name string, onItem func(UnlockItem)) ([]UnlockItem, error) {
		out := []UnlockItem{
			{Name: "Disney+", Status: "No", CheckTime: "2026-11-02 10:00:00"},
			{Name: "disney+", Status: "Pending"},
			{Name: "Netflix", Status: "Yes", Region: "US", CheckTime: "2026-11-02 10:00:01"},
		}
		if onItem != nil {
			onItem(out[2])
		}
		return out, nil
	}
	if err := tt.Click("unlock-test-all"); err != nil {
		t.Fatal(err)
	}
	if len(m.Unlock) != 2 || m.Unlock[0].Name != "Disney+" || m.Unlock[0].Status != "No" {
		t.Fatalf("merged all %+v", m.Unlock)
	}
	if m.Unlock[1].Status != "Yes" || m.Notice != "" || tt.HasText("测试中...") {
		t.Fatalf("done %+v notice %q\n%s", m.Unlock, m.Notice, join(tt.Texts()))
	}
}

func TestUnlockEmptyDisablesTestAll(t *testing.T) {
	m := NewModel(demoHome())
	m.Hooks.RefreshUnlock = func() { m.Unlock = nil }
	calls := 0
	m.Hooks.CheckUnlock = func(string, func(UnlockItem)) ([]UnlockItem, error) {
		calls++
		return nil, nil
	}
	m.Router.Replace("/unlock")
	tt := NewPageTester(m)
	if err := tt.Click("unlock-test-all"); err != nil {
		t.Fatal(err)
	}
	if calls != 0 {
		t.Fatal("empty list started a test")
	}
}

func TestUnlockReloadOnEnterNotOnShow(t *testing.T) {
	m := NewModel(demoHome())
	n := 0
	m.Hooks.RefreshUnlock = func() {
		n++
		m.Unlock = []UnlockItem{{Name: "Netflix", Status: "Pending"}}
	}
	tt := NewPageTester(m)
	if n != 0 {
		t.Fatalf("home loaded unlock %d times", n)
	}
	if err := tt.Click("测 试"); err != nil {
		t.Fatal(err)
	}
	if n != 1 || m.Router.Path() != "/unlock" {
		t.Fatalf("enter %d path %s", n, m.Router.Path())
	}
	m.PageVisible = false
	if err := tt.Click("解锁测试"); err != nil {
		t.Fatal(err)
	}
	m.PageVisible = true
	if err := tt.Click("解锁测试"); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("visibility reloaded unlock (%d)", n)
	}
	if err := tt.Click("首页"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("测 试"); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("re-enter loaded %d times, want 2", n)
	}
}

func TestUnlockCacheRoundTrip(t *testing.T) {
	prefs := DefaultPrefs()
	prefs.UnlockResults = []UnlockItem{{
		Name: "netflix", Status: "Yes", Region: "US", CheckTime: "2026-11-02 08:00:00",
	}}
	raw, err := json.Marshal(prefs)
	if err != nil {
		t.Fatal(err)
	}
	loaded := DefaultPrefs()
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatal(err)
	}
	m := NewModel(demoHome())
	m.Prefs = loaded
	m.Router.Replace("/unlock")
	tt := NewPageTester(m)
	if !tt.HasText("Netflix") || !tt.HasText("支持") || !tt.HasText("US") || !tt.HasText("-- --") {
		t.Fatalf("restored\n%s", join(tt.Texts()))
	}
	for _, item := range m.Unlock {
		if item.Name == "Netflix" && item.Status == "Yes" && item.Region == "US" {
			return
		}
	}
	t.Fatalf("model %+v", m.Unlock)
}

func sameRow(a, b float32) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 4
}

func namesOf(items []UnlockItem) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.Name
	}
	return out
}
