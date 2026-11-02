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
		if err := tt.Click(p.Title); err != nil {
			t.Fatal(err)
		}
		if m.Router.Path() != p.Path {
			t.Fatalf("path %q after clicking %s, want %s", m.Router.Path(), p.Title, p.Path)
		}
		if !tt.HasText(p.Body) {
			t.Fatalf("page %s missing %q\ntexts: %q", p.ID, p.Body, tt.Texts())
		}
	}
	if err := tt.Click("首页"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("1.19.33") {
		t.Fatal("home version disappeared")
	}
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
