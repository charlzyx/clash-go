package profile

import (
	"bytes"
	"strings"
)

// ClashPatch is the slice of settings that must land in the runtime
// config.yaml the kernel reloads. Zero-value ports that are enabled are
// written as 0, which mihomo treats as "listener off".
//
// Redir and TProxy are applied only when the matching Apply flag is set,
// so a platform that does not show the control does not wipe a key the
// profile itself defined.
type ClashPatch struct {
	AllowLan     bool
	IPv6         bool
	UnifiedDelay bool
	// LogLevel is the settings-page value. "warning" and the legacy
	// spelling "warn" both become the kernel spelling. mihomo v1.19
	// UnmarshalText accepts "warning" (LogLevel.String) and rejects "warn".
	LogLevel    string
	MixedPort   int
	SocksPort   int
	SocksOn     bool
	HTTPPort    int
	HTTPOn      bool
	RedirPort   int
	RedirOn     bool
	ApplyRedir  bool
	TProxyPort  int
	TProxyOn    bool
	ApplyTProxy bool
	// DNSOverwrite replaces the top-level dns block with the app's
	// per-profile overwrite. When false the profile's own dns block is
	// left in place.
	DNSOverwrite bool
}

// dnsOverwriteBlock is the runtime DNS used when a profile turns overwrite
// on. listen is omitted so the kernel does not bind port 53. Nameserver is
// required whenever dns.enable is true.
const dnsOverwriteBlock = `dns:
  enable: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  fake-ip-filter:
    - "*.lan"
    - "*.local"
    - localhost.ptlogin2.qq.com
  nameserver:
    - 223.6.6.6
    - 8.8.8.8`

// KernelLogLevel maps a settings-page log level onto the spelling this
// mihomo accepts. Unknown values fall back to info.
func KernelLogLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug", "info", "error", "silent":
		return strings.ToLower(strings.TrimSpace(level))
	case "warn", "warning":
		return "warning"
	default:
		return "info"
	}
}

// ApplyClash rewrites top-level clash keys on a copy of data. Other keys,
// comments, and indented blocks stay. The profile file on disk is not this
// function's input when the caller still has the original subscription.
func ApplyClash(data []byte, p ClashPatch) []byte {
	lines, trailing := splitLines(data)
	lines = setTopScalar(lines, "allow-lan", yamlBool(p.AllowLan))
	lines = setTopScalar(lines, "ipv6", yamlBool(p.IPv6))
	lines = setTopScalar(lines, "unified-delay", yamlBool(p.UnifiedDelay))
	lines = setTopScalar(lines, "log-level", KernelLogLevel(p.LogLevel))
	mixed := p.MixedPort
	if mixed < 1 || mixed > 65535 {
		mixed = 7897
	}
	lines = setTopScalar(lines, "mixed-port", itoa(mixed))
	lines = setTopScalar(lines, "socks-port", itoa(enabledPort(p.SocksOn, p.SocksPort)))
	lines = setTopScalar(lines, "port", itoa(enabledPort(p.HTTPOn, p.HTTPPort)))
	if p.ApplyRedir {
		lines = setTopScalar(lines, "redir-port", itoa(enabledPort(p.RedirOn, p.RedirPort)))
	}
	if p.ApplyTProxy {
		lines = setTopScalar(lines, "tproxy-port", itoa(enabledPort(p.TProxyOn, p.TProxyPort)))
	}
	if p.DNSOverwrite {
		lines = setTopBlock(lines, "dns", dnsOverwriteBlock)
	}
	return joinLines(lines, trailing || len(lines) > 0)
}

// HasTopLevelDNS reports whether data has a top-level dns key. Indented
// dns keys (inside a proxy, for example) do not count. The settings page
// asks before overwrite when this is true.
func HasTopLevelDNS(data []byte) bool {
	lines, _ := splitLines(data)
	_, _, ok := topKeyRange(lines, "dns")
	return ok
}

func enabledPort(on bool, port int) int {
	if !on || port < 1 || port > 65535 {
		return 0
	}
	return port
}

func yamlBool(on bool) string {
	if on {
		return "true"
	}
	return "false"
}

func splitLines(data []byte) (lines [][]byte, trailing bool) {
	if len(data) == 0 {
		return nil, false
	}
	trailing = data[len(data)-1] == '\n'
	s := strings.TrimSuffix(string(data), "\n")
	if s == "" {
		return nil, trailing
	}
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, []byte(line))
	}
	return lines, trailing
}

func joinLines(lines [][]byte, trailing bool) []byte {
	b := bytes.Join(lines, []byte("\n"))
	if trailing {
		b = append(b, '\n')
	}
	return b
}

func setTopScalar(lines [][]byte, key, value string) [][]byte {
	line := []byte(key + ": " + value)
	if start, end, ok := topKeyRange(lines, key); ok {
		out := append([][]byte{}, lines[:start]...)
		out = append(out, line)
		return append(out, lines[end:]...)
	}
	return append(lines, line)
}

func setTopBlock(lines [][]byte, key, block string) [][]byte {
	blockLines, _ := splitLines([]byte(block))
	if start, end, ok := topKeyRange(lines, key); ok {
		out := append([][]byte{}, lines[:start]...)
		out = append(out, blockLines...)
		return append(out, lines[end:]...)
	}
	return append(append([][]byte{}, lines...), blockLines...)
}

func topKeyRange(lines [][]byte, key string) (start, end int, ok bool) {
	for i, line := range lines {
		if indented(line) || blank(line) {
			continue
		}
		rest, found := splitTopKey(line, key)
		if !found {
			continue
		}
		if rest != "" && !strings.HasPrefix(rest, "|") && !strings.HasPrefix(rest, ">") {
			return i, i + 1, true
		}
		j := i + 1
		for j < len(lines) {
			if blank(lines[j]) {
				k := j
				for k < len(lines) && blank(lines[k]) {
					k++
				}
				if k < len(lines) && indented(lines[k]) {
					j = k
					continue
				}
				break
			}
			if !indented(lines[j]) {
				break
			}
			j++
		}
		return i, j, true
	}
	return 0, 0, false
}

func splitTopKey(line []byte, key string) (rest string, ok bool) {
	trim := bytes.TrimSpace(line)
	prefix := key + ":"
	if !bytes.HasPrefix(trim, []byte(prefix)) {
		return "", false
	}
	return string(bytes.TrimSpace(trim[len(prefix):])), true
}

func indented(line []byte) bool {
	return len(line) > 0 && (line[0] == ' ' || line[0] == '\t')
}

func blank(line []byte) bool {
	return len(bytes.TrimSpace(line)) == 0
}
