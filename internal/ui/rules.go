package ui

import (
	"sort"

	mygoui "github.com/egoist/mygo/ui"
)

// Sizes taken from the rules outline and rule-item.tsx. mygo's list
// estimates a row from the height it measures, so each row is at least
// RuleRowEstimate; Verge's VirtualList estimateSize is 40.
const (
	RuleSearchRowHeight float32 = 36
	RuleSearchPadX      float32 = 10
	RuleRowEstimate     float32 = 40
	RuleIndexMinWidth   float32 = 30
	RulePayloadMinWidth float32 = 120
	RuleProviderListMin float32 = 250
	RuleScrollTopAfter  float32 = 100
)

// syncRuleRefresh loads rules and rule providers when the page is entered
// and when it is shown again after the window was hidden.
//
// Verge's effect calls both refreshers on every run, and again inside the
// same run when the page is already visible, so a mount while visible
// refreshes twice. One refresh covers that enter. A later hide then show
// refreshes again.
func (m *Model) syncRuleRefresh() {
	if m == nil || m.Router == nil {
		return
	}
	on := m.Router.Path() == "/rules"
	if on && (!m.rulesOnPage || (m.PageVisible && !m.rulesPageVisible)) {
		m.refreshRules()
	}
	m.rulesOnPage = on
	if on {
		m.rulesPageVisible = m.PageVisible
	} else {
		m.rulesPageVisible = false
	}
}

func (m *Model) refreshRules() {
	if m.Hooks.RefreshRules != nil {
		m.Hooks.RefreshRules()
	}
}

func (m *Model) rulesPage(c *mygoui.Context, pal Palette) {
	// The list scrolls itself. The shared page() helper would put another
	// scroll container around it.
	mygoui.Column(c).Fill().Children(func() {
		head := mygoui.Row(c).Height(58).Padding(0, 20).Gap(8).AlignItems(mygoui.Center).Shrink(0)
		head.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
		head.Children(func() {
			mygoui.Text(c, "规则").FontSize(PageTitleSize).Bold().Shrink(0)
			mygoui.Spacer(c).Grow(1)
			// Verge hides the button when there is no rule provider.
			if len(m.RuleProviders) > 0 && textButton(c, pal, "规则集合", "规则集合", "rule-providers") {
				m.Dialogs.RuleProvider = true
			}
		})
		mygoui.Row(c).Height(RuleSearchRowHeight).Padding(0, RuleSearchPadX).AlignItems(mygoui.Center).Shrink(0).
			Label("rules-search-row").Children(func() {
			mygoui.TextInput(c, &m.RuleQuery).FillWidth().Placeholder("过滤条件").Label("rules-search")
		})
		rows := FilterRules(m.Rules, m.RuleQuery)
		if len(rows) == 0 {
			// Verge renders BaseEmpty (inbox + 「空空如也」). The outline says
			// the region stays empty, and this window has no shared empty
			// component, so a search with no payload match leaves the region
			// blank. The search row stays.
			mygoui.Column(c).Grow(1).FillWidth().Label("rules-empty")
			return
		}
		mygoui.Column(c).Grow(1).FillWidth().Children(func() {
			m.RulesList.Key = func(i int) any {
				if i < 0 || i >= len(rows) {
					return i
				}
				return rows[i].Line
			}
			first, _ := m.RulesList.Visible()
			mygoui.List(c, &m.RulesList, len(rows), func(i int) {
				m.ruleItem(c, pal, rows[i])
			}).Grow(1).FillWidth().Label("rules-list")
			if ShowRuleScrollTop(first) {
				b := mygoui.ButtonBase(c.Key("rules-scroll-top")).Absolute().Bottom(20).Right(20).
					Size(36, 36).Radius(18).Label("rules-scroll-top").
					Background(pal.Text.Alpha(0.1)).TextColor(pal.Text).
					AlignItems(mygoui.Center).Justify(mygoui.Center).Tooltip("回到顶部")
				b.Children(func() { mygoui.Text(c, "↑").FontSize(18) })
				if b.Clicked() {
					m.RulesList.ScrollTo(0, mygoui.Start)
					c.Invalidate()
				}
			}
		})
	})
}

func (m *Model) ruleItem(c *mygoui.Context, pal Palette, row RuleLine) {
	payload := row.Rule.Payload
	if payload == "" {
		payload = "-"
	}
	mygoui.Row(c).MinHeight(RuleRowEstimate).Padding(4, 16).Gap(18).AlignItems(mygoui.Center).
		BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider).Children(func() {
		mygoui.Text(c, itoa(row.Line)).FontSize(13).TextColor(pal.TextMuted).
			MinWidth(RuleIndexMinWidth).TextAlign(mygoui.Center)
		mygoui.Column(c).Grow(1).MinWidth(0).Children(func() {
			mygoui.Text(c, payload).FontSize(14).MinWidth(RulePayloadMinWidth).SingleLine()
			mygoui.Row(c).Gap(24).AlignItems(mygoui.Center).Children(func() {
				mygoui.Text(c, row.Rule.Type).FontSize(13).TextColor(pal.TextMuted).MinWidth(RulePayloadMinWidth).SingleLine()
				mygoui.Text(c, row.Rule.Proxy).FontSize(13).TextColor(RuleProxyColor(pal, row.Rule.Proxy)).SingleLine()
			})
		})
	})
}

func (m *Model) ruleProviderDialog(c *mygoui.Context, pal Palette) {
	if !m.Dialogs.RuleProvider {
		return
	}
	mygoui.Modal(c, &m.Dialogs.RuleProvider, func() {
		mygoui.Column(c).Width(480).Padding(16).Gap(8).Background(pal.Dialog).Radius(8).Children(func() {
			mygoui.Row(c).AlignItems(mygoui.Center).Gap(8).Children(func() {
				mygoui.Text(c, "规则集合").FontSize(18).Bold().Grow(1)
				if filledButton(c, pal, "更新全部", "rule-prov-all", true) {
					m.updateAllRuleProviders()
					c.Invalidate()
				}
			})
			providers := sortedRuleProviders(m.RuleProviders)
			mygoui.Column(c).MinHeight(RuleProviderListMin).Gap(8).FillWidth().Label("rule-provider-list").Children(func() {
				if len(providers) == 0 {
					mygoui.Text(c, "没有可更新的规则集合").TextColor(pal.TextMuted)
				}
				for i := range providers {
					p := providers[i]
					mygoui.Row(c).FillWidth().Padding(8, 12).Gap(8).Radius(8).Background(pal.Node).
						AlignItems(mygoui.Center).Children(func() {
						mygoui.Column(c).Grow(1).Gap(4).Children(func() {
							mygoui.Row(c).Gap(8).AlignItems(mygoui.Center).Children(func() {
								mygoui.Text(c, p.Name).FontSize(14).SingleLine()
								ruleBadge(c, pal, itoa(p.RuleCount))
								mygoui.Spacer(c).Grow(1)
								mygoui.Text(c, "更新于 "+FormatUpdated(p.UpdatedAt, c.Now())).FontSize(12).TextColor(pal.TextMuted)
							})
							mygoui.Row(c).Gap(4).Children(func() {
								ruleBadge(c, pal, p.VehicleType)
								ruleBadge(c, pal, p.Behavior)
							})
						})
						label := "更新"
						if p.Updating {
							label = "…"
						}
						if textButton(c, pal, label, "更新", "rule-prov-update-"+p.Name) {
							m.updateOneRuleProvider(p.Name)
							c.Invalidate()
						}
					})
				}
			})
			if filledButton(c, pal, "关闭", "rule-provider-close", false) {
				m.Dialogs.RuleProvider = false
			}
		})
	})
}

func ruleBadge(c *mygoui.Context, pal Palette, text string) {
	mygoui.Row(c).Padding(0, 2).Radius(4).Shrink(0).Border(1, pal.Secondary.Alpha(0.5)).Children(func() {
		mygoui.Text(c, text).FontSize(10).TextColor(pal.Secondary.Alpha(0.8))
	})
}

func sortedRuleProviders(in []RuleProvider) []RuleProvider {
	out := append([]RuleProvider(nil), in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (m *Model) updateOneRuleProvider(name string) {
	if m.Hooks.UpdateRuleProvider != nil {
		if err := m.Hooks.UpdateRuleProvider(name); err != nil {
			m.SetNotice(name + " 更新失败: " + err.Error())
			return
		}
	}
	m.SetNotice(name + " 更新成功")
	m.refreshRules()
}

func (m *Model) updateAllRuleProviders() {
	if len(m.RuleProviders) == 0 {
		m.SetNotice("没有可更新的规则集合")
		return
	}
	// Verge logs a single provider's error and still reports that every
	// provider was updated once the loop finishes.
	names := make([]string, len(m.RuleProviders))
	for i, p := range m.RuleProviders {
		names[i] = p.Name
	}
	for _, name := range names {
		if m.Hooks.UpdateRuleProvider == nil {
			continue
		}
		_ = m.Hooks.UpdateRuleProvider(name)
	}
	m.SetNotice("所有规则集合均已更新")
	m.refreshRules()
}
