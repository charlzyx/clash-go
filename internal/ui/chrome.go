package ui

import (
	"fmt"

	mygoui "github.com/egoist/mygo/ui"
)

func (m *Model) page(c *mygoui.Context, pal Palette, title string, actions func(), body func()) {
	mygoui.Column(c).Fill().Children(func() {
		head := mygoui.Row(c).Height(58).Padding(0, 20).Gap(8).AlignItems(mygoui.Center).Shrink(0)
		head.BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider)
		head.Children(func() {
			mygoui.Text(c, title).FontSize(PageTitleSize).Bold().Shrink(0)
			mygoui.Spacer(c).Grow(1)
			if actions != nil {
				actions()
			}
		})
		if m.Notice != "" {
			mygoui.Text(c, m.Notice).FontSize(13).TextColor(pal.Error).Padding(6, 20)
		}
		if m.Home.Error != "" && m.Router.Path() == "/home" {
			mygoui.Text(c, m.Home.Error).FontSize(13).TextColor(pal.Error).Padding(0, 20, 6, 20)
		}
		mygoui.Scroll(c).Grow(1).Children(func() {
			if body != nil {
				body()
			}
		})
	})
}

func textButton(c *mygoui.Context, pal Palette, label, tip, id string) bool {
	b := mygoui.ButtonBase(c.Key(id)).Height(28).Padding(0, 8).Radius(6).Label(id)
	b.Background(mygoui.Transparent).TextColor(pal.Text)
	b.Tooltip(tip)
	b.Children(func() {
		mygoui.Text(c, label).FontSize(13)
	})
	return b.Clicked()
}

func filledButton(c *mygoui.Context, pal Palette, label, id string, on bool) bool {
	b := mygoui.ButtonBase(c.Key(id)).Height(28).Padding(0, 10).Radius(6).Label(id)
	if on {
		b.Background(pal.Primary).TextColor(pal.OnPrimary)
	} else {
		b.Background(mygoui.Transparent).TextColor(pal.Text).Border(1, pal.Primary)
	}
	b.Children(func() { mygoui.Text(c, label).FontSize(13) })
	return b.Clicked()
}

func card(c *mygoui.Context, pal Palette, title, mark string, tone Tone, action func(), body func()) mygoui.Element {
	col := mygoui.Column(c).FillWidth().Radius(8).Background(pal.Card)
	col.Children(func() {
		mygoui.Row(c).Padding(8, 16).Gap(12).AlignItems(mygoui.Center).
			BorderWidth(0, 0, 1, 0).BorderColor(pal.Divider).Children(func() {
			mygoui.Column(c).Size(38, 38).Radius(6).Shrink(0).Background(pal.Color(tone).Alpha(0.12)).
				AlignItems(mygoui.Center).Justify(mygoui.Center).Children(func() {
				mygoui.Text(c, mark).FontSize(13).TextColor(pal.Color(tone)).Bold()
			})
			mygoui.Text(c, title).FontSize(18).FontWeight(500).SingleLine().Grow(1)
			if action != nil {
				action()
			}
		})
		mygoui.Column(c).Padding(16).Gap(8).AlignItems(mygoui.Start).Children(func() {
			if body != nil {
				body()
			}
		})
	})
	return col
}

func chip(c *mygoui.Context, pal Palette, text string, tone Tone, fill bool) {
	bg := pal.Color(tone).Alpha(0.12)
	fg := pal.Color(tone)
	if fill {
		bg = pal.Color(tone)
		fg = pal.OnPrimary
	}
	mygoui.Row(c).Padding(2, 8).Radius(4).Background(bg).Shrink(0).Children(func() {
		mygoui.Text(c, text).FontSize(12).TextColor(fg)
	})
}

func meter(c *mygoui.Context, pal Palette, pct int) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	mygoui.Row(c).FillWidth().Height(8).Radius(4).Background(pal.Primary.Alpha(0.12)).Children(func() {
		mygoui.Box(c).Height(8).WidthPercent(float32(pct)).Background(pal.Primary).Radius(4)
	})
	mygoui.Text(c, fmt.Sprintf("%d%%", pct)).FontSize(12).TextColor(pal.TextMuted)
}

func signalMark(c *mygoui.Context, pal Palette, ms int) {
	tone := pal.Color(SignalTone(ms))
	bars := SignalBars(ms)
	mygoui.Row(c).Size(38, 38).Radius(6).Shrink(0).Gap(2).Padding(6).
		Background(tone.Alpha(0.12)).AlignItems(mygoui.End).Justify(mygoui.Center).
		Tooltip(SignalLabel(ms) + ": " + DelayText(ms, DelayTimeout)).Children(func() {
		for i := 0; i < 4; i++ {
			h := float32(8 + i*4)
			bar := mygoui.Box(c).Width(3).Height(h).Radius(1)
			if i < bars {
				bar.Background(tone)
			} else {
				bar.Background(pal.TextMuted.Alpha(0.35))
			}
		}
	})
}

func modeFace(c *mygoui.Context, pal Palette, label, id, current string, wide, badge bool) bool {
	on := current == id
	var clicked bool
	mygoui.Column(c).Grow(1).MaxWidth(160).AlignItems(mygoui.Center).Children(func() {
		b := mygoui.ButtonBase(c.Key(id)).FillWidth().Height(36).Radius(6).Label(id)
		if on {
			b.Background(pal.Primary).TextColor(pal.OnPrimary)
		} else {
			b.Background(pal.Card).TextColor(pal.Text).Border(1, pal.Divider)
		}
		b.AlignItems(mygoui.Center).Justify(mygoui.Center)
		b.Children(func() {
			weight := 400
			if on {
				weight = 600
			}
			mygoui.Text(c, label).FontSize(14).FontWeight(weight)
			if badge {
				mygoui.Box(c).Size(8, 8).Radius(4).Background(pal.Success)
			}
		})
		clicked = b.Clicked()
		mark := mygoui.Box(c).Margin(4, 0, 0, 0).Height(2).Background(pal.Primary)
		if wide {
			mark.Width(16)
		} else {
			mark.Width(9)
		}
		if !on {
			mark.Background(mygoui.Transparent)
		}
	})
	return clicked
}

func dot(c *mygoui.Context, pal Palette, on bool) {
	if !on {
		return
	}
	mygoui.Box(c).Size(8, 8).Radius(4).Background(pal.Success).Shrink(0)
}

func sectionCard(c *mygoui.Context, pal Palette, title string, body func()) mygoui.Element {
	col := mygoui.Column(c).FillWidth().Radius(8).Background(pal.Card).Padding(12, 16).Gap(2)
	col.Children(func() {
		mygoui.Text(c, title).FontSize(16).Bold().Padding(4, 0, 8, 0)
		if body != nil {
			body()
		}
	})
	return col
}

// clickRow is a settings row. click nil means the row is not a button and
// trail draws the control. A non-nil click shows a chevron, or "…" while busy.
func clickRow(c *mygoui.Context, pal Palette, label, id string, busy bool, click func(), trail func()) {
	if click == nil {
		mygoui.Row(c).FillWidth().Padding(5, 0).AlignItems(mygoui.Center).Gap(8).Children(func() {
			mygoui.Text(c, label).FontSize(14).Grow(1)
			if trail != nil {
				trail()
			}
		})
		return
	}
	mygoui.Row(c).FillWidth().AlignItems(mygoui.Center).Gap(4).Children(func() {
		b := mygoui.ButtonBase(c.Key(id)).Grow(1).Padding(5, 0).Disabled(busy).Label(id)
		b.Children(func() {
			mygoui.Text(c, label).FontSize(14).Grow(1)
			if trail != nil {
				trail()
			}
			if busy {
				mygoui.Text(c, "…").FontSize(14).TextColor(pal.Primary)
			} else {
				mygoui.Text(c, "›").FontSize(18).TextColor(pal.TextMuted)
			}
		})
		if b.Clicked() {
			click()
		}
	})
}
