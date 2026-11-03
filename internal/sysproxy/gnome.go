package sysproxy

import (
	"fmt"
	"os"
	"strings"
)

// Gnome sets the proxy through GNOME gsettings. It is the Linux controller.
// Other desktops (KDE, Sway, a bare tty) are refused: writing the GNOME
// keys there would report success without changing the proxy those sessions
// actually use.
type Gnome struct {
	Run Runner
	// Desktop overrides the session name. DesktopFromEnv reads
	// XDG_CURRENT_DESKTOP, then XDG_SESSION_DESKTOP, when Desktop is empty.
	Desktop        string
	DesktopFromEnv bool
}

// NewGnome returns the controller Platform uses on Linux.
func NewGnome() *Gnome { return &Gnome{Run: execRunner, DesktopFromEnv: true} }

// Current reads the GNOME manual proxy. A session that is not GNOME reports
// the proxy as off.
func (g *Gnome) Current() (Settings, error) {
	if g == nil || g.Run == nil {
		return Settings{}, fmt.Errorf("sysproxy: no command runner")
	}
	if !gnomeDesktop(g.desktop()) {
		return Settings{}, nil
	}
	mode, err := g.get("org.gnome.system.proxy", "mode")
	if err != nil {
		return Settings{}, err
	}
	httpHost, err := g.get("org.gnome.system.proxy.http", "host")
	if err != nil {
		return Settings{}, err
	}
	httpPort, err := g.get("org.gnome.system.proxy.http", "port")
	if err != nil {
		return Settings{}, err
	}
	socksHost, err := g.get("org.gnome.system.proxy.socks", "host")
	if err != nil {
		return Settings{}, err
	}
	socksPort, err := g.get("org.gnome.system.proxy.socks", "port")
	if err != nil {
		return Settings{}, err
	}
	cur := Settings{
		Enabled: gvariant(mode) == "manual",
		HTTP:    joinHostPort(gvariant(httpHost), gvariant(httpPort)),
		SOCKS:   joinHostPort(gvariant(socksHost), gvariant(socksPort)),
	}
	return cur, nil
}

// Set writes the GNOME manual proxy, or turns it off. A non-GNOME session
// returns an error so the switch can roll back instead of pretending.
func (g *Gnome) Set(next Settings) error {
	if g == nil || g.Run == nil {
		return fmt.Errorf("sysproxy: no command runner")
	}
	if !gnomeDesktop(g.desktop()) {
		return fmt.Errorf("当前桌面不支持系统代理（仅 GNOME）")
	}
	if !next.Enabled {
		return g.set("org.gnome.system.proxy", "mode", "none")
	}
	httpHost, httpPort, err := splitEndpoint(next.HTTP)
	if err != nil {
		return fmt.Errorf("sysproxy: http endpoint: %w", err)
	}
	socksHost, socksPort, err := splitEndpoint(next.SOCKS)
	if err != nil {
		return fmt.Errorf("sysproxy: socks endpoint: %w", err)
	}
	steps := []struct{ schema, key, value string }{
		{"org.gnome.system.proxy", "use-same-proxy", "false"},
		{"org.gnome.system.proxy.http", "host", httpHost},
		{"org.gnome.system.proxy.http", "port", httpPort},
		{"org.gnome.system.proxy.https", "host", httpHost},
		{"org.gnome.system.proxy.https", "port", httpPort},
		{"org.gnome.system.proxy.socks", "host", socksHost},
		{"org.gnome.system.proxy.socks", "port", socksPort},
		{"org.gnome.system.proxy", "ignore-hosts", "['localhost', '127.0.0.0/8', '::1', '*.local']"},
		{"org.gnome.system.proxy", "mode", "manual"},
	}
	for _, step := range steps {
		if err := g.set(step.schema, step.key, step.value); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gnome) desktop() string {
	if g.Desktop != "" {
		return g.Desktop
	}
	if !g.DesktopFromEnv {
		return ""
	}
	if d := os.Getenv("XDG_CURRENT_DESKTOP"); d != "" {
		return d
	}
	return os.Getenv("XDG_SESSION_DESKTOP")
}

func (g *Gnome) get(schema, key string) (string, error) {
	out, err := g.Run("gsettings", "get", schema, key)
	if err != nil {
		return "", fmt.Errorf("读取系统代理失败: %w", err)
	}
	return out, nil
}

func (g *Gnome) set(schema, key, value string) error {
	if _, err := g.Run("gsettings", "set", schema, key, value); err != nil {
		return fmt.Errorf("设置系统代理失败: %w", err)
	}
	return nil
}

func gnomeDesktop(desktop string) bool {
	d := strings.ToLower(desktop)
	return strings.Contains(d, "gnome") || strings.Contains(d, "unity")
}

func gvariant(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "uint32 ")
	s = strings.Trim(s, "'\"")
	return strings.TrimSpace(s)
}

func joinHostPort(host, port string) string {
	host, port = strings.TrimSpace(host), strings.TrimSpace(port)
	if host == "" || port == "" || port == "0" {
		return ""
	}
	return host + ":" + port
}
