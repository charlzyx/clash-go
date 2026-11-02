package app

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
	"github.com/charlzyx/clash-go/internal/ui"
)

// TestUnlockResultsPersistAcrossRestart stores the cache in settings.json
// prefs (the clash_verge_unlock_results stand-in) and reads it back on a
// second model. The kernel is not started.
func TestUnlockResultsPersistAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	store := appconfig.NewMemory(dir)
	prefs := ui.DefaultPrefs()
	prefs.UnlockResults = []ui.UnlockItem{
		{Name: "Netflix", Status: "Yes", Region: "US", CheckTime: "2026-11-02 08:00:00"},
		{Name: "Disney+", Status: "Risky", CheckTime: "2026-11-02 08:00:01"},
	}
	raw, err := json.Marshal(prefs)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(appconfig.Settings{Prefs: string(raw)}); err != nil {
		t.Fatal(err)
	}

	first := newUnlockHarness(t, dir, store)
	m := first.Model()
	m.Router.Replace("/unlock")
	tt := ui.NewPageTester(m)
	if !tt.HasText("支持") || !tt.HasText("存疑") || !tt.HasText("Netflix") {
		t.Fatalf("seeded cache\n%s", strings.Join(tt.Texts(), "\n"))
	}
	m.Hooks.CheckUnlock = func(name string, onItem func(ui.UnlockItem)) ([]ui.UnlockItem, error) {
		item := ui.UnlockItem{Name: name, Status: "No", CheckTime: "2026-11-02 09:00:00"}
		if onItem != nil {
			onItem(item)
		}
		return []ui.UnlockItem{item}, nil
	}
	if err := tt.Click("unlock-Netflix"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("不支持") {
		t.Fatalf("retest\n%s", strings.Join(tt.Texts(), "\n"))
	}
	saved, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(saved.Prefs, "clash_verge_unlock_results") || !strings.Contains(saved.Prefs, `"No"`) {
		t.Fatalf("settings prefs %s", saved.Prefs)
	}

	second := newUnlockHarness(t, dir, store)
	again := second.Model()
	again.Router.Replace("/unlock")
	tt2 := ui.NewPageTester(again)
	var netflix, disney ui.UnlockItem
	for _, item := range again.Unlock {
		switch item.Name {
		case "Netflix":
			netflix = item
		case "Disney+":
			disney = item
		}
	}
	if netflix.Status != "No" || disney.Status != "Risky" {
		t.Fatalf("restart netflix %+v disney %+v\n%s", netflix, disney, strings.Join(tt2.Texts(), "\n"))
	}
}

func newUnlockHarness(t *testing.T, dir string, store appconfig.Store) *App {
	t.Helper()
	a := New(Options{
		Dirs:  appdir.Fixed(dir),
		Store: store,
		Proxy: sysproxy.NewMemory(),
		TUN:   tun.NewMemory(),
		Core:  core.New(),
	})
	a.dir = dir
	settings, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	a.settings = settings
	a.model = ui.NewModel(ui.HomeState{ConfigDir: dir})
	a.applySettings()
	a.bindHooks()
	return a
}
