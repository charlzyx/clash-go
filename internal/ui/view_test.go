package ui

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	mygoui "github.com/egoist/mygo/ui"
)

// P0 ui.NewTester cases for ProfilesView and ProxiesView are deferred.
// docs/test/priority.md asks for an import field, a proxy row click, and a
// delay result. /proxies and /profile still render PlaceholderPage, so those
// interactions do not exist yet and this package does not invent the views.

func demoHome() HomeState {
	return HomeState{
		Version:     "1.19.33",
		Running:     true,
		MixedPort:   7897,
		Mode:        "rule",
		ConfigDir:   "~/Library/Application Support/app.clashgo.desktop",
		SystemProxy: false,
		TUN:         false,
	}
}

func TestHomeShowsKernelStatus(t *testing.T) {
	m := NewModel(demoHome())
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	for _, text := range []string{"首页", "内核版本", "1.19.33", "运行状态", "运行中", "7897", "rule"} {
		if !tt.HasText(text) {
			t.Errorf("home missing %q\ntexts: %q", text, tt.Texts())
		}
	}
}

func TestNavigatePages(t *testing.T) {
	m := NewModel(demoHome())
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	for _, p := range Pages() {
		if p.ID == "home" {
			continue
		}
		if err := tt.Click(p.Nav); err != nil {
			t.Fatal(err)
		}
		if m.Router.Path() != p.Path {
			t.Fatalf("path %q after clicking %s, want %s", m.Router.Path(), p.Nav, p.Path)
		}
		if !tt.HasText(p.Title) || !tt.HasText(p.Body) {
			t.Fatalf("page %s missing %q or %q\ntexts: %q", p.ID, p.Title, p.Body, tt.Texts())
		}
	}
	if err := tt.Click("首页"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("1.19.33") {
		t.Fatal("home version disappeared")
	}
}

// TestM1ShellLocks freezes the three M1 shell fixes from issue #5.
// ui.NewTester does not expose a text node's font size, so the 20px
// token is the exported PageTitleSize used by the headings.
func TestM1ShellLocks(t *testing.T) {
	t.Run("subscriptions path", func(t *testing.T) {
		var subs Page
		var found bool
		for _, p := range Pages() {
			if p.Path == "/subscriptions" {
				t.Fatal("subscriptions nav must use /profile, not /subscriptions")
			}
			if p.ID == "subscriptions" || p.Nav == "订阅" {
				subs = p
				found = true
			}
		}
		if !found || subs.Path != "/profile" || subs.Nav != "订阅" {
			t.Fatalf("subscriptions page %+v, want path /profile and nav 订阅", subs)
		}
		m := NewModel(demoHome())
		tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
		if err := tt.Click("订阅"); err != nil {
			t.Fatal(err)
		}
		if m.Router.Path() != "/profile" {
			t.Fatalf("path %q after clicking 订阅, want /profile", m.Router.Path())
		}
		if !hasExactText(tt, "订阅") {
			t.Fatalf("texts %q", tt.Texts())
		}
	})

	t.Run("proxies title", func(t *testing.T) {
		var title string
		for _, p := range Pages() {
			if p.ID == "proxies" {
				title = p.Title
			}
		}
		if title != "代理组" {
			t.Fatalf("proxies title %q, want 代理组", title)
		}
		m := NewModel(demoHome())
		tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
		if err := tt.Click("代理"); err != nil {
			t.Fatal(err)
		}
		if m.Router.Path() != "/proxies" {
			t.Fatalf("path %q", m.Router.Path())
		}
		if !hasExactText(tt, "代理组") {
			t.Fatalf("texts %q, want an exact 代理组 title", tt.Texts())
		}
	})

	t.Run("page title size", func(t *testing.T) {
		if PageTitleSize != 20 {
			t.Fatalf("PageTitleSize %v, want 20", PageTitleSize)
		}
		m := NewModel(demoHome())
		tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
		if !hasExactText(tt, "首页") {
			t.Fatalf("home title missing, texts %q", tt.Texts())
		}
		if err := tt.Click("代理"); err != nil {
			t.Fatal(err)
		}
		if !hasExactText(tt, "代理组") {
			t.Fatalf("proxies title missing, texts %q", tt.Texts())
		}
	})
}

func hasExactText(tt *mygoui.Tester, s string) bool {
	for _, text := range tt.Texts() {
		if text == s {
			return true
		}
	}
	return false
}

func TestStoppedHome(t *testing.T) {
	m := NewModel(HomeState{Version: "1.19.33", Error: "内核没有启动"})
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	if !tt.HasText("已停止") || !tt.HasText("内核没有启动") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

// TestPageScreenshots renders every page at the Verge window size, in
// light and dark. Set CLASHGO_SCREENSHOT_DIR to also write
// <id>-light.png and <id>-dark.png.
func TestPageScreenshots(t *testing.T) {
	dir := os.Getenv("CLASHGO_SCREENSHOT_DIR")
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var light, dark *image.RGBA
	for _, theme := range []struct {
		name string
		dark bool
	}{{"light", false}, {"dark", true}} {
		for _, p := range Pages() {
			img := renderPage(demoHome(), p.Path, theme.dark)
			b := img.Bounds()
			if b.Dx() != WindowWidth || b.Dy() != WindowHeight {
				t.Fatalf("%s-%s is %dx%d, want %dx%d", p.ID, theme.name, b.Dx(), b.Dy(), WindowWidth, WindowHeight)
			}
			if p.ID == "home" && !theme.dark {
				light = img
			}
			if p.ID == "home" && theme.dark {
				dark = img
			}
			if dir == "" {
				continue
			}
			path := filepath.Join(dir, p.ID+"-"+theme.name+".png")
			if err := writePNG(path, img); err != nil {
				t.Fatal(err)
			}
		}
	}
	if light == nil || dark == nil {
		t.Fatal("missing home frames")
	}
	if sameRGBA(light, dark) {
		t.Fatal("light and dark home frames are identical")
	}
}

func renderPage(home HomeState, path string, dark bool) *image.RGBA {
	m := NewModel(home)
	if path != "/home" {
		m.Router.Replace(path)
	}
	tt := mygoui.NewTester(m.View, WindowWidth, WindowHeight)
	tt.SetScale(1)
	tt.SetDark(dark)
	for i := 0; i < 5 && tt.FramePending(); i++ {
		tt.Advance(time.Second)
	}
	return tt.Image()
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func sameRGBA(a, b *image.RGBA) bool {
	if a == nil || b == nil || !a.Rect.Eq(b.Rect) || len(a.Pix) != len(b.Pix) {
		return false
	}
	for i := range a.Pix {
		if a.Pix[i] != b.Pix[i] {
			return false
		}
	}
	return true
}
