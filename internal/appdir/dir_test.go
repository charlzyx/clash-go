package appdir

import "testing"

func TestPath(t *testing.T) {
	home := "/Users/me"
	if got, want := Path(home, "darwin"), "/Users/me/Library/Application Support/app.clashgo.desktop"; got != want {
		t.Fatalf("darwin path %q, want %q", got, want)
	}
	if got, want := Path(home, "linux"), "/Users/me/.config/clash-go"; got != want {
		t.Fatalf("linux path %q, want %q", got, want)
	}
	if got, want := Path(home, "windows"), "/Users/me/AppData/Roaming/app.clashgo.desktop"; got != want {
		t.Fatalf("windows path %q, want %q", got, want)
	}
}

func TestFixed(t *testing.T) {
	if _, err := Fixed("").ConfigDir(); err == nil {
		t.Fatal("empty Fixed returned a directory")
	}
	got, err := Fixed("/tmp/clash-go").ConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != "/tmp/clash-go" {
		t.Fatalf("got %q", got)
	}
}
