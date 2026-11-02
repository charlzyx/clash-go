package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/profile"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
	"github.com/charlzyx/clash-go/internal/ui"

	mygoui "github.com/egoist/mygo/ui"
)

// TestProfileImportShowsInListAndSurvivesRestart drives /profile, imports a
// subscription URL, and reloads a new model from the same profiles.json.
// The kernel is not started. Do not call t.Parallel: mihomo is a process-wide
// singleton, and this package must not race it.
func TestProfileImportShowsInListAndSurvivesRestart(t *testing.T) {
	dir := t.TempDir()
	body := []byte("mode: rule\nmixed-port: 7897\nproxies: []\n")
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "text/yaml")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	raw := srv.URL + "/sub"
	name := ui.URLHost(raw)

	first := newProfileHarness(dir)
	m := first.Model()
	tt := openProfiles(t, m)
	if listed, _ := profileListed(tt, m, name, raw); listed {
		t.Fatalf("list already shows %q before import", name)
	}
	if err := tt.Click("profile-url"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("profile-url") {
		t.Fatal("subscription URL field did not take focus")
	}
	tt.Type(raw)
	if m.ProfileURL != raw {
		t.Fatalf("url field %q, want %q", m.ProfileURL, raw)
	}
	if err := tt.Click("do-import"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/sub" {
		t.Fatalf("import fetched %q, want /sub", gotPath)
	}
	if m.ProfileURL != "" {
		t.Fatalf("url field still %q after import, notice %q", m.ProfileURL, m.Notice)
	}
	if ok, why := profileListed(tt, m, name, raw); !ok {
		t.Fatalf("imported %q missing from the list (%s), notice %q, texts %q, profiles %+v", name, why, m.Notice, tt.Texts(), m.Profiles)
	}
	if len(m.Profiles) != 1 {
		t.Fatalf("profiles %+v, want one new entry", m.Profiles)
	}

	saved, err := os.ReadFile(filepath.Join(dir, "profiles.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(saved), raw) || !strings.Contains(string(saved), name) {
		t.Fatalf("profiles.json missing the import:\n%s", saved)
	}

	// Restart: a new model loaded only from the same profiles.json.
	restored := newProfileHarness(dir)
	restored.loadProfiles()
	again := restored.Model()
	tt2 := openProfiles(t, again)
	if ok, why := profileListed(tt2, again, name, raw); !ok {
		t.Fatalf("after restart %q is missing (%s), texts %q, profiles %+v", name, why, tt2.Texts(), again.Profiles)
	}
	if len(again.Profiles) != 1 || again.Profiles[0].URL != raw || again.CurrentProfile == "" {
		t.Fatalf("restored %+v current %q", again.Profiles, again.CurrentProfile)
	}
	if again.Home.Profile == nil || again.Home.Profile.Name != name {
		t.Fatalf("home profile %+v", again.Home.Profile)
	}
}

// newProfileHarness wires the window to a profiles.json store without
// Prepare, so the kernel stays stopped.
func newProfileHarness(dir string) *App {
	a := New(Options{
		Dirs:  appdir.Fixed(dir),
		Store: appconfig.NewMemory(dir),
		Proxy: sysproxy.NewMemory(),
		TUN:   tun.NewMemory(),
		Core:  core.New(),
	})
	a.dir = dir
	a.profiles = profile.NewStore(dir)
	a.model = ui.NewModel(ui.HomeState{ConfigDir: dir})
	a.bindHooks()
	return a
}

func openProfiles(t *testing.T, m *ui.Model) *mygoui.Tester {
	t.Helper()
	tt := ui.NewPageTester(m)
	if m.Router.Path() != "/profile" {
		if err := tt.Click("订阅"); err != nil {
			t.Fatal(err)
		}
	}
	if m.Router.Path() != "/profile" {
		t.Fatalf("path %q, want /profile", m.Router.Path())
	}
	if _, ok := tt.Find("profile-url"); !ok || !tt.HasText("导入") || !tt.HasText("新建") {
		t.Fatalf("profile import row missing, texts %q", tt.Texts())
	}
	return tt
}

func profileListed(tt *mygoui.Tester, m *ui.Model, name, raw string) (bool, string) {
	if m == nil || m.Router.Path() != "/profile" {
		return false, "not on /profile"
	}
	if _, ok := tt.Find("profile-" + name); !ok {
		return false, "no profile card"
	}
	if !tt.HasText(name) {
		return false, "card name is not drawn"
	}
	for _, p := range m.Profiles {
		if p.Name == name && p.Type == "remote" && p.URL == raw {
			return true, ""
		}
	}
	return false, "model has no remote profile for that URL"
}
