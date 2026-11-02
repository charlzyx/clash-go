package profile

import (
	"errors"
	"strings"
	"testing"
)

func TestImportRemote(t *testing.T) {
	s := NewStore(t.TempDir())
	if _, err := s.ImportRemote("not a url", func(string) ([]byte, error) { return nil, nil }); !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("err %v", err)
	}
	item, err := s.ImportRemote("https://example.com/sub", func(string) ([]byte, error) {
		return []byte("proxies: []\n"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.Type != "remote" || item.Name != "example.com" {
		t.Fatalf("%+v", item)
	}
	f, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if f.Current != item.UID || len(f.Items) != 1 {
		t.Fatalf("%+v", f)
	}
	if _, err := s.ImportLocal("notes.txt", []byte("x")); err == nil || !strings.Contains(err.Error(), "YAML") {
		t.Fatalf("txt err %v", err)
	}
	local, err := s.ImportLocal("/tmp/home.yaml", []byte("mode: rule\n"))
	if err != nil {
		t.Fatal(err)
	}
	if local.Type != "local" || local.Name != "home.yaml" {
		t.Fatalf("%+v", local)
	}
}

func TestParseAndMode(t *testing.T) {
	raw := []byte(`
mode: rule
mixed-port: 7897
proxies:
  - name: HK
    type: ss
proxy-groups:
  - name: Proxy
    type: select
    now: HK
    proxies:
      - HK
      - DIRECT
rules:
  - MATCH,Proxy
  - DOMAIN,example.com,DIRECT
`)
	got, err := ParseClashYAML(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != "rule" || got.MixedPort != 7897 || got.Rules != 2 || len(got.Groups) != 1 || got.Groups[0].Now != "HK" {
		t.Fatalf("%+v", got)
	}
	if got.Groups[0].Proxies[0].Type != "ss" || got.Groups[0].Proxies[1].Name != "DIRECT" {
		t.Fatalf("%+v", got.Groups[0].Proxies)
	}
	next := SetMode(raw, "global")
	again, err := ParseClashYAML(next)
	if err != nil || again.Mode != "global" {
		t.Fatalf("mode %q %v", again.Mode, err)
	}
	if !strings.Contains(string(EnsureMixedPort([]byte("mode: rule\n"), 7897)), "mixed-port: 7897") {
		t.Fatal("port not prepended")
	}
}
