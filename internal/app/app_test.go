package app

import (
	"errors"
	"testing"

	"github.com/charlzyx/clash-go/internal/appconfig"
	"github.com/charlzyx/clash-go/internal/appdir"
	"github.com/charlzyx/clash-go/internal/core"
	"github.com/charlzyx/clash-go/internal/sysproxy"
	"github.com/charlzyx/clash-go/internal/tun"
)

type fakeProxy struct {
	cur    sysproxy.Settings
	setErr error
}

func (f *fakeProxy) Current() (sysproxy.Settings, error) { return f.cur, nil }

func (f *fakeProxy) Set(s sysproxy.Settings) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.cur = s
	return nil
}

type fakeTUN struct {
	on     bool
	setErr error
}

func (f *fakeTUN) Current() (bool, error) { return f.on, nil }

func (f *fakeTUN) Set(on bool) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.on = on
	return nil
}

func TestSetSystemProxyUsesController(t *testing.T) {
	dir := t.TempDir()
	proxy := &fakeProxy{}
	store := appconfig.NewMemory(dir)
	a := New(Options{
		Dirs:  appdir.Fixed(dir),
		Store: store,
		Proxy: proxy,
		TUN:   tun.NewMemory(),
		Core:  core.New(),
	})
	if err := a.SetSystemProxy(true); err != nil {
		t.Fatal(err)
	}
	if !proxy.cur.Enabled || proxy.cur.HTTP != "127.0.0.1:7897" || proxy.cur.SOCKS != "127.0.0.1:7897" {
		t.Fatalf("proxy %+v", proxy.cur)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !got.SystemProxy || got.TUN {
		t.Fatalf("settings %+v", got)
	}

	proxy.setErr = errors.New("denied")
	if err := a.SetSystemProxy(false); err == nil {
		t.Fatal("SetSystemProxy succeeded after the controller failed")
	}
	got, _ = store.Load()
	if !got.SystemProxy {
		t.Fatal("failed Set cleared the saved proxy")
	}
}

func TestSetTUNRollsBack(t *testing.T) {
	dir := t.TempDir()
	store := appconfig.NewMemory(dir)
	store.FailSave(errors.New("disk"))
	tunCtl := &fakeTUN{}
	a := New(Options{
		Dirs:  appdir.Fixed(dir),
		Store: store,
		Proxy: sysproxy.NewMemory(),
		TUN:   tunCtl,
		Core:  core.New(),
	})
	if err := a.SetTUN(true); err == nil {
		t.Fatal("SetTUN succeeded with a failing store")
	}
	if tunCtl.on {
		t.Fatal("TUN stayed on after the store failed")
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.TUN {
		t.Fatal("store kept TUN")
	}
}
