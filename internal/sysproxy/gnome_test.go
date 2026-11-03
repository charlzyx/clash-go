package sysproxy

import (
	"errors"
	"runtime"
	"strings"
	"testing"
)

func TestPlatformSelectsRealController(t *testing.T) {
	c := Platform()
	switch runtime.GOOS {
	case "darwin":
		if _, ok := c.(*System); !ok {
			t.Fatalf("darwin platform %T", c)
		}
	case "linux":
		if _, ok := c.(*Gnome); !ok {
			t.Fatalf("linux platform %T", c)
		}
	default:
		if _, ok := c.(Unsupported); !ok {
			t.Fatalf("platform %T", c)
		}
	}
}

func TestGnomeSetAndCurrent(t *testing.T) {
	var calls []string
	g := &Gnome{Desktop: "ubuntu:GNOME", Run: func(name string, args ...string) (string, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		if len(args) > 0 && args[0] == "get" {
			switch args[2] {
			case "mode":
				return "'manual'\n", nil
			case "host":
				return "'127.0.0.1'\n", nil
			case "port":
				return "uint32 7897\n", nil
			}
		}
		return "", nil
	}}
	if err := g.Set(Settings{Enabled: true, HTTP: "127.0.0.1:7897", SOCKS: "127.0.0.1:7897"}); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(calls, "\n")
	for _, want := range []string{
		"gsettings set org.gnome.system.proxy.http host 127.0.0.1",
		"gsettings set org.gnome.system.proxy.http port 7897",
		"gsettings set org.gnome.system.proxy.https host 127.0.0.1",
		"gsettings set org.gnome.system.proxy.socks host 127.0.0.1",
		"gsettings set org.gnome.system.proxy.socks port 7897",
		"gsettings set org.gnome.system.proxy mode manual",
		"gsettings set org.gnome.system.proxy ignore-hosts ['localhost', '127.0.0.0/8', '::1', '*.local']",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q\n%s", want, joined)
		}
	}
	cur, err := g.Current()
	if err != nil {
		t.Fatal(err)
	}
	if !cur.Enabled || cur.HTTP != "127.0.0.1:7897" || cur.SOCKS != "127.0.0.1:7897" {
		t.Fatalf("current %+v", cur)
	}

	calls = nil
	if err := g.Set(Settings{}); err != nil {
		t.Fatal(err)
	}
	if strings.Join(calls, "\n") != "gsettings set org.gnome.system.proxy mode none" {
		t.Fatalf("disable calls %q", calls)
	}
}

func TestGnomeRefusesOtherDesktops(t *testing.T) {
	g := &Gnome{Desktop: "KDE", Run: func(string, ...string) (string, error) {
		t.Fatal("gsettings ran on KDE")
		return "", nil
	}}
	if err := g.Set(Settings{Enabled: true, HTTP: "127.0.0.1:7897", SOCKS: "127.0.0.1:7897"}); err == nil || !strings.Contains(err.Error(), "GNOME") {
		t.Fatalf("err %v", err)
	}
	cur, err := g.Current()
	if err != nil || cur.Enabled {
		t.Fatalf("current %+v err %v", cur, err)
	}

	bare := &Gnome{Run: func(string, ...string) (string, error) {
		return "", errors.New("no gsettings")
	}}
	if err := bare.Set(Settings{Enabled: true, HTTP: "127.0.0.1:1", SOCKS: "127.0.0.1:1"}); err == nil {
		t.Fatal("empty desktop succeeded")
	}
}

func TestUnsupportedSetFails(t *testing.T) {
	u := Unsupported{GOOS: "windows"}
	if err := u.Set(Settings{Enabled: true}); err == nil || !strings.Contains(err.Error(), "windows") {
		t.Fatalf("err %v", err)
	}
	cur, err := u.Current()
	if err != nil || cur.Enabled {
		t.Fatalf("current %+v %v", cur, err)
	}
}
