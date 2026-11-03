package profile

import (
	"strings"
	"testing"
)

func TestKernelLogLevel(t *testing.T) {
	if KernelLogLevel("warning") != "warning" || KernelLogLevel("warn") != "warning" {
		t.Fatalf("warning map %q %q", KernelLogLevel("warning"), KernelLogLevel("warn"))
	}
	if KernelLogLevel("DEBUG") != "debug" || KernelLogLevel("nope") != "info" {
		t.Fatalf("other levels %q %q", KernelLogLevel("DEBUG"), KernelLogLevel("nope"))
	}
}

func TestApplyClashRewritesTopLevelOnly(t *testing.T) {
	raw := []byte(`# keep this comment
mode: global
mixed-port: 1
socks-port: 7898
port: 7899
redir-port: 7895
tproxy-port: 7896
allow-lan: false
ipv6: true
unified-delay: true
log-level: info
proxies:
  - name: hk
    type: ss
    port: 443
    allow-lan: false
rules:
  - MATCH,DIRECT
dns:
  enable: false
  nameserver:
    - 1.1.1.1
`)
	got := string(ApplyClash(raw, ClashPatch{
		AllowLan: true, IPv6: false, UnifiedDelay: false, LogLevel: "warn",
		MixedPort: 17897, SocksOn: true, SocksPort: 11080, HTTPOn: false, HTTPPort: 7899,
		ApplyRedir: false, ApplyTProxy: false,
		DNSOverwrite: false,
	}))
	for _, want := range []string{
		"# keep this comment",
		"mode: global",
		"allow-lan: true",
		"ipv6: false",
		"unified-delay: false",
		"log-level: warning",
		"mixed-port: 17897",
		"socks-port: 11080",
		"port: 0",
		"redir-port: 7895",
		"tproxy-port: 7896",
		"nameserver:",
		"- 1.1.1.1",
		"type: ss",
		"    port: 443",
		"    allow-lan: false",
		"MATCH,DIRECT",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\n%s", want, got)
		}
	}
	if strings.Contains(got, "log-level: warn\n") || strings.Contains(got, "log-level: info") {
		t.Fatalf("log level was not normalized\n%s", got)
	}
}

func TestApplyClashPortsAndDNS(t *testing.T) {
	raw := []byte("mode: rule\nredir-port: 9\n")
	got := string(ApplyClash(raw, ClashPatch{
		LogLevel: "error", MixedPort: 0,
		SocksOn: false, HTTPOn: true, HTTPPort: 8080,
		ApplyRedir: true, RedirOn: false, RedirPort: 9,
		ApplyTProxy: true, TProxyOn: true, TProxyPort: 7896,
		DNSOverwrite: true,
	}))
	for _, want := range []string{
		"mixed-port: 7897",
		"socks-port: 0",
		"port: 8080",
		"redir-port: 0",
		"tproxy-port: 7896",
		"log-level: error",
		"dns:",
		"enable: true",
		"enhanced-mode: fake-ip",
		"fake-ip-range: 198.18.0.1/16",
		"223.6.6.6",
		"8.8.8.8",
		"mode: rule",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\n%s", want, got)
		}
	}
	if strings.Contains(got, "nameserver:\n    - 1.1.1.1") {
		t.Fatal("overwrite kept the old nameserver")
	}
	if !HasTopLevelDNS([]byte(got)) {
		t.Fatal("overwrite did not leave a top-level dns key")
	}
	if HasTopLevelDNS([]byte("proxies:\n  - name: a\n    dns: 1.1.1.1\n")) {
		t.Fatal("indented dns counted as top-level")
	}
	if !HasTopLevelDNS([]byte("dns: {enable: false}\n")) {
		t.Fatal("flow dns key missed")
	}
}
