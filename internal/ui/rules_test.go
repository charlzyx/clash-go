package ui

import (
	"errors"
	"fmt"
	"testing"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

func TestRulesPageListsAndFiltersPayload(t *testing.T) {
	m := NewModel(demoHome())
	m.Rules = []Rule{
		{Type: "Domain", Payload: "Google.Example", Proxy: "Proxy"},
		{Type: "DomainSuffix", Payload: "Yahoo.Example", Proxy: "DIRECT"},
		{Type: "Match", Payload: "", Proxy: "DIRECT"},
	}
	m.Router.Replace("/rules")
	tt := NewPageTester(m)
	for _, text := range []string{"规则", "Google.Example", "Yahoo.Example", "Domain", "DomainSuffix", "Match", "Proxy", "DIRECT", "-"} {
		if !tt.HasText(text) {
			t.Errorf("missing %q\n%s", text, join(tt.Texts()))
		}
	}
	if tt.HasText("分流规则会显示在这里") {
		t.Fatal("rules page is still the placeholder")
	}
	for _, line := range []string{"1", "2", "3"} {
		if !hasExactText(tt, line) {
			t.Fatalf("missing line %s\n%s", line, join(tt.Texts()))
		}
	}
	row, ok := tt.Find("rules-search-row")
	if !ok || row.H < RuleSearchRowHeight-2 || row.H > RuleSearchRowHeight+4 {
		t.Fatalf("search row %#v", row)
	}
	if dash, ok := tt.Find("-"); !ok || dash.W < RulePayloadMinWidth-1 {
		t.Fatalf("payload min width %#v", dash)
	}
	if idx, ok := tt.Find("1"); !ok || idx.W < RuleIndexMinWidth-1 {
		t.Fatalf("index min width %#v", idx)
	}
	if err := tt.Click("rules-search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("yahoo")
	if tt.HasText("Google.Example") || !tt.HasText("Yahoo.Example") {
		t.Fatalf("filter\n%s", join(tt.Texts()))
	}
	if !hasExactText(tt, "2") || hasExactText(tt, "1") {
		t.Fatalf("line number changed with the filter\n%s", join(tt.Texts()))
	}

	m.RuleQuery = "missing.payload"
	if err := tt.Click("规则"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Google.Example") || tt.HasText("Yahoo.Example") {
		t.Fatalf("empty match still lists rules\n%s", join(tt.Texts()))
	}
	if _, ok := tt.Find("rules-empty"); !ok {
		t.Fatal("empty match should leave a blank region")
	}
	if tt.HasText("空空如也") {
		t.Fatal("empty match invented an empty-state label")
	}
	if _, ok := tt.Find("rules-search"); !ok {
		t.Fatal("search box disappeared with the empty match")
	}
}

func TestRulesProviderUpdate(t *testing.T) {
	m := NewModel(demoHome())
	m.RuleProviders = []RuleProvider{{
		Name: "reject", VehicleType: "File", Behavior: "Domain", RuleCount: 4,
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}}
	var got []string
	failed := false
	m.Hooks.UpdateRuleProvider = func(name string) error {
		got = append(got, name)
		if failed {
			return errors.New("boom")
		}
		return nil
	}
	m.Router.Replace("/rules")
	tt := NewPageTester(m)
	if _, ok := tt.Find("rule-providers"); !ok {
		t.Fatal("provider button missing")
	}
	if err := tt.Click("rule-providers"); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"规则集合", "reject", "File", "Domain", "更新全部", "更新于"} {
		if !tt.HasText(text) {
			t.Errorf("dialog missing %q\n%s", text, join(tt.Texts()))
		}
	}
	list, ok := tt.Find("rule-provider-list")
	if !ok || list.H < RuleProviderListMin-1 {
		t.Fatalf("provider list %#v, want min height %v", list, RuleProviderListMin)
	}
	if err := tt.Click("rule-prov-update-reject"); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "reject" || !tt.HasText("reject 更新成功") {
		t.Fatalf("single update got %v\n%s", got, join(tt.Texts()))
	}
	if err := tt.Click("rule-prov-all"); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1] != "reject" || !tt.HasText("所有规则集合均已更新") {
		t.Fatalf("update all got %v\n%s", got, join(tt.Texts()))
	}

	failed = true
	if err := tt.Click("rule-prov-update-reject"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("reject 更新失败: boom") {
		t.Fatalf("failure notice\n%s", join(tt.Texts()))
	}
}

func TestRulesRefreshOnEnterAndShow(t *testing.T) {
	m := NewModel(demoHome())
	n := 0
	m.Hooks.RefreshRules = func() {
		n++
		m.Rules = []Rule{{Type: "Domain", Payload: "refreshed.example", Proxy: "DIRECT"}}
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if n != 0 {
		t.Fatalf("home page refreshed rules %d times", n)
	}
	if err := tt.Click("规则"); err != nil {
		t.Fatal(err)
	}
	if n != 1 || !tt.HasText("refreshed.example") {
		t.Fatalf("enter refresh %d\n%s", n, join(tt.Texts()))
	}
	m.PageVisible = false
	if err := tt.Click("规则"); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("hiding the window refreshed rules (%d)", n)
	}
	m.PageVisible = true
	if err := tt.Click("规则"); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("showing the window again refreshed %d times, want 2", n)
	}
}

func TestRulesScrollTop(t *testing.T) {
	m := NewModel(demoHome())
	for i := 0; i < 40; i++ {
		m.Rules = append(m.Rules, Rule{
			Type: "Domain", Payload: fmt.Sprintf("host-%02d.example", i), Proxy: "DIRECT",
		})
	}
	m.Router.Replace("/rules")
	tt := NewPageTester(m)
	if _, ok := tt.Find("rules-scroll-top"); ok {
		t.Fatal("scroll-to-top is visible before scrolling")
	}
	list, ok := tt.Find("rules-list")
	if !ok {
		t.Fatal("missing rules list")
	}
	tt.Scroll(list.X+20, list.Y+30, 0, 500)
	if _, ok := tt.Find("rules-scroll-top"); !ok {
		first, last := m.RulesList.Visible()
		t.Fatalf("scroll-to-top missing after 500px, visible rows %d-%d", first, last)
	}
	if err := tt.Click("rules-scroll-top"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("rules-scroll-top"); ok {
		first, _ := m.RulesList.Visible()
		t.Fatalf("scroll-to-top stayed after the click, first row %d", first)
	}
	if !tt.HasText("host-00.example") {
		t.Fatal("list did not return to the first rule")
	}
}

func TestRulesHiddenWithoutProviders(t *testing.T) {
	m := NewModel(demoHome())
	m.Router.Replace("/rules")
	tt := NewPageTester(m)
	if _, ok := tt.Find("rule-providers"); ok || tt.HasText("规则集合") {
		t.Fatal("provider button shows with no providers")
	}
}
