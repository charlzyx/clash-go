package ui

import (
	"runtime"

	mygoui "github.com/egoist/mygo/ui"
)

func (m *Model) settingsPage(c *mygoui.Context, pal Palette) {
	m.page(c, pal, "设置", func() {
		if iconButton(c, pal, iconHelp, "使用手册", "set-manual") {
			c.OpenURL("https://clash-verge-rev.github.io/index.html")
		}
		if iconButton(c, pal, iconTelegram, "Telegram 频道", "set-tg") {
			c.OpenURL("https://t.me/clash_verge_re")
		}
		if iconButton(c, pal, iconGitHub, "GitHub 项目地址", "set-gh") {
			c.OpenURL("https://github.com/clash-verge-rev/clash-verge-rev")
		}
	}, func() {
		w, _ := c.Size()
		cols := 1
		if WideLayout(w) {
			cols = 2
		}
		mygoui.Grid(c).Columns(cols).Gap(12).Padding(16).Children(func() {
			mygoui.Column(c).Gap(12).Children(func() {
				m.systemSettings(c, pal)
				m.clashSettings(c, pal)
			})
			mygoui.Column(c).Gap(12).Children(func() {
				m.vergeBasic(c, pal)
				m.vergeAdvanced(c, pal)
			})
		})
	})
}

func (m *Model) systemSettings(c *mygoui.Context, pal Palette) {
	sectionCard(c, pal, "系统设置", func() {
		m.switchLine(c, pal, "虚拟网卡模式", &m.Home.TUN, "tun-row", func(on bool) { m.setTUN(on) }, func() {
			if iconButton(c, pal, settingsGear(), "TUN", "tun-gear") {
				m.Dialogs.TUN = true
			}
		})
		m.switchLine(c, pal, "系统代理", &m.Home.SystemProxy, "sysproxy-row", func(on bool) { m.setProxy(on) }, func() {
			if iconButton(c, pal, settingsGear(), "系统代理", "sysproxy-gear") {
				m.Dialogs.Sysproxy = true
			}
		})
		m.switchLine(c, pal, "开机自启", &m.Prefs.AutoLaunch, "launch-row", func(bool) { m.savePrefs() }, nil)
		m.switchLine(c, pal, "静默启动", &m.Prefs.SilentStart, "silent-row", func(bool) { m.savePrefs() }, func() {
			hintIcon(c, pal, iconInfo, "程序启动时以后台模式运行，不显示程序面板")
		})
	})
}

func (m *Model) switchLine(c *mygoui.Context, pal Palette, label string, on *bool, id string, apply func(bool), extra func()) {
	mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Gap(8).Children(func() {
		mygoui.Text(c, label).FontSize(14).Grow(1)
		if extra != nil {
			extra()
		}
		sw := mygoui.Switch(c, on).Label(id)
		if sw.Changed() && apply != nil {
			apply(*on)
		}
	})
	_ = pal
}

func (m *Model) clashSettings(c *mygoui.Context, pal Palette) {
	sectionCard(c, pal, "Clash 设置", func() {
		m.switchLine(c, pal, "局域网连接", &m.Prefs.AllowLan, "allow-lan", func(bool) { m.savePrefs() }, func() {
			if iconButton(c, pal, iconLan, "网络接口", "lan-gear") {
				m.SetNotice("网络接口")
			}
		})
		disabled := m.currentProfile() == nil || m.DNSAsk
		dnsOn := m.dnsOverwriteTarget()
		mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Gap(8).Children(func() {
			mygoui.Text(c, "DNS 覆写").FontSize(14).Grow(1)
			hintIcon(c, pal, iconInfo, "DNS 覆写开关按订阅独立保存，仅对对应订阅生效。")
			if iconButton(c, pal, settingsGear(), "DNS 覆写", "dns-gear") {
				m.Dialogs.DNS = true
			}
			sw := mygoui.Switch(c, dnsOn).Disabled(disabled).Label("dns-overwrite")
			if sw.Changed() {
				if *dnsOn && m.ProfileHasDNS {
					*dnsOn = false
					m.Prefs.DNSOverwrite = false
					m.DNSAsk = true
					return
				}
				m.Prefs.DNSOverwrite = *dnsOn
				m.savePrefs()
			}
		})
		m.switchLine(c, pal, "IPv6", &m.Prefs.IPv6, "ipv6", func(bool) { m.savePrefs() }, nil)
		m.switchLine(c, pal, "统一延迟", &m.Prefs.UnifiedDelay, "unified", func(bool) { m.savePrefs() }, func() {
			hintIcon(c, pal, iconInfo, "开启统一延迟时，会进行两次延迟测试，以消除连接握手等带来的不同类型节点的延迟差异")
		})
		m.labeledSelect(c, "日志等级", &m.Prefs.LogLevel, []string{"debug", "info", "warning", "error", "silent"}, "log-level")
		clickRow(c, pal, "端口设置", "ports", m.Busy == "ports", func() { m.Dialogs.Ports = true }, func() {
			mygoui.Text(c, itoa(m.Prefs.MixedPort)).FontSize(13).TextColor(pal.TextMuted)
		})
		clickRow(c, pal, "外部控制", "external", false, func() { m.SetNotice("外部控制器监听地址") }, func() {
			if iconButton(c, pal, settingsGear(), "外部控制跨域设置", "cors") {
				m.SetNotice("外部控制跨域设置")
			}
		})
		clickRow(c, pal, "网页界面", "webui", false, func() { m.SetNotice("网页界面") }, nil)
		clickRow(c, pal, "Clash 内核", "core-row", false, func() { m.Dialogs.Lite = false }, func() {
			mygoui.Text(c, orDash(m.Home.Version)).FontSize(13)
			if iconButton(c, pal, settingsGear(), "Clash 内核", "core-gear") {
				m.SetNotice("Clash 内核 " + orDash(m.Home.Version))
			}
		})
		if runtime.GOOS == "windows" {
			clickRow(c, pal, "UWP 工具", "uwp", false, func() {}, nil)
		}
		clickRow(c, pal, "更新 GeoData", "geo", m.Busy == "geo", func() {
			m.SetNotice("GeoData 已更新")
		}, nil)
		clickRow(c, pal, "流量隧道管理", "tunnels", false, func() { m.SetNotice("流量隧道管理") }, nil)
	})
}

func (m *Model) vergeBasic(c *mygoui.Context, pal Palette) {
	sectionCard(c, pal, "Verge 基础设置", func() {
		names := make([]string, len(Languages()))
		for i, l := range Languages() {
			names[i] = l.Name
		}
		cur := LanguageName(m.Prefs.Language)
		mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Children(func() {
			mygoui.Text(c, "语言设置").FontSize(14).Grow(1)
			sel := mygoui.Select(c, &cur, names).Width(140).Label("language")
			if sel.Changed() {
				m.Prefs.Language = LanguageID(cur)
				m.savePrefs()
			}
		})
		mygoui.Text(c, "主题模式").FontSize(14)
		mygoui.Row(c).Gap(6).Padding(4, 0).Children(func() {
			for _, item := range []struct{ id, label string }{{"light", "浅色"}, {"dark", "深色"}, {"system", "系统"}} {
				if filledButton(c, pal, item.label, "theme-"+item.id, m.Prefs.ThemeMode == item.id) {
					m.Prefs.ThemeMode = item.id
					m.savePrefs()
				}
			}
		})
		if runtime.GOOS != "linux" || m.ForceTray {
			m.choiceRow(c, "托盘点击事件", &m.Prefs.TrayEvent, TrayEvents(), "tray-event")
		}
		envs := CopyEnvTypes()
		mygoui.Row(c).FillWidth().AlignItems(mygoui.Center).Padding(5, 0).Children(func() {
			mygoui.Text(c, "复制环境变量类型").FontSize(14).Grow(1)
			sel := mygoui.Select(c, &m.Prefs.CopyEnv, envs).Width(140).Label("copy-env")
			if sel.Changed() {
				m.savePrefs()
			}
			if iconButton(c, pal, iconCopy, "复制环境变量", "copy-env-btn") {
				m.copy("export https_proxy=http://" + m.proxyAddr())
			}
		})
		m.choiceRow(c, "启动页面", &m.Prefs.StartPage, StartPages(), "start-page")
		mygoui.Row(c).FillWidth().AlignItems(mygoui.Center).Gap(6).Padding(5, 0).Children(func() {
			mygoui.Text(c, "启动脚本").FontSize(14)
			mygoui.TextInput(c, &m.Prefs.StartupScript).Grow(1).ReadOnly(true).Label("startup-script")
			if filledButton(c, pal, "浏览", "browse-script", false) {
				if m.Hooks.BrowseScript != nil {
					path, err := m.Hooks.BrowseScript()
					if err != nil {
						m.SetNotice(err.Error())
					} else {
						m.Prefs.StartupScript = path
						m.savePrefs()
					}
				}
			}
			if m.Prefs.StartupScript != "" && textButton(c, pal, "清除", "清除", "clear-script") {
				m.Prefs.StartupScript = ""
				m.savePrefs()
			}
		})
		clickRow(c, pal, "主题设置", "theme-row", false, func() { m.Dialogs.Theme = true }, nil)
		clickRow(c, pal, "界面设置", "layout-row", false, func() { m.Dialogs.Layout = true }, nil)
		clickRow(c, pal, "杂项设置", "misc-row", false, func() { m.Dialogs.Misc = true }, nil)
		clickRow(c, pal, "热键设置", "hotkey-row", false, func() { m.Dialogs.Hotkey = true }, nil)
	})
}

func (m *Model) vergeAdvanced(c *mygoui.Context, pal Palette) {
	sectionCard(c, pal, "Verge 高级设置", func() {
		clickRow(c, pal, "备份设置", "backup", false, func() { m.Dialogs.Backup = true }, func() {
			hintIcon(c, pal, iconInfo, "支持本地或 WebDAV 方式备份配置文件")
		})
		clickRow(c, pal, "当前配置", "runtime-row", false, func() { m.Dialogs.Runtime = true }, nil)
		clickRow(c, pal, "配置目录", "conf-dir", false, func() { m.openDir("config") }, func() {
			hintIcon(c, pal, iconInfo, "如果软件运行异常，备份并删除此文件夹下的所有文件，重启软件")
		})
		clickRow(c, pal, "内核目录", "core-dir", false, func() { m.openDir("core") }, nil)
		clickRow(c, pal, "日志目录", "log-dir", false, func() { m.openDir("logs") }, nil)
		clickRow(c, pal, "检查更新", "check-update", m.Busy == "check-update", func() {
			m.SetNotice("当前已是最新版本")
		}, nil)
		clickRow(c, pal, "开发者工具", "devtools", false, func() { m.SetNotice("开发者工具") }, nil)
		clickRow(c, pal, "轻量模式设置", "lite-row", false, func() { m.Dialogs.Lite = true }, func() {
			hintIcon(c, pal, iconInfo, "关闭 GUI 界面，仅保留内核运行")
		})
		clickRow(c, pal, "退出", "quit-row", false, func() {
			if m.Hooks.Quit != nil {
				m.Hooks.Quit()
			}
			m.SetNotice("退出")
		}, nil)
		clickRow(c, pal, "导出诊断信息", "diag", false, nil, func() {
			if iconButton(c, pal, iconCopy, "导出诊断信息", "copy-diag") {
				m.copy("Clash Go " + AppVersion)
			}
		})
		ver := "v" + trimV(m.Home.VergeVersion)
		clickRow(c, pal, "Verge 版本", "verge-version", false, nil, func() {
			mygoui.Text(c, ver).FontSize(14)
			if iconButton(c, pal, iconCopy, "复制 Verge 版本号", "copy-version") {
				m.copy(ver)
				m.SetNotice("Verge 版本已复制到剪贴板")
			}
		})
	})
}

// dnsOverwriteTarget is the current subscription's switch. With no
// subscription the control is disabled and writes a scratch bool.
func (m *Model) dnsOverwriteTarget() *bool {
	if cur := m.currentProfile(); cur != nil {
		return &cur.Option.DNSOverwrite
	}
	return &m.Prefs.DNSOverwrite
}

func (m *Model) openDir(kind string) {
	if m.Hooks.OpenDir != nil {
		m.Hooks.OpenDir(kind)
	}
	m.SetNotice("已打开目录")
}

func (m *Model) labeledSelect(c *mygoui.Context, label string, value *string, options []string, id string) {
	mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Children(func() {
		mygoui.Text(c, label).FontSize(14).Grow(1)
		sel := mygoui.Select(c, value, options).Width(120).Label(id)
		if sel.Changed() {
			m.savePrefs()
		}
	})
}

func (m *Model) choiceRow(c *mygoui.Context, label string, id *string, choices []LanguageChoice, key string) {
	names := make([]string, len(choices))
	cur := ""
	for i, ch := range choices {
		names[i] = ch.Name
		if ch.ID == *id {
			cur = ch.Name
		}
	}
	if cur == "" && len(names) > 0 {
		cur = names[0]
	}
	mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Children(func() {
		mygoui.Text(c, label).FontSize(14).Grow(1)
		sel := mygoui.Select(c, &cur, names).Width(140).Label(key)
		if sel.Changed() {
			for _, ch := range choices {
				if ch.Name == cur {
					*id = ch.ID
				}
			}
			m.savePrefs()
		}
	})
}
