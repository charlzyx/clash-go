package ui

import (
	"testing"

	mygoui "github.com/egoist/mygo/ui"
)

// TestSidebarLogoSizes locks the sidebar brand slot from
// docs/spec/outlines/layout.md and design-tokens.md §8.1 (#8).
//
// Expanded: the favicon is 36×36, then 「Clash Go」 about 5px to its right.
// Collapsed: NewCollapsedTester's rail is 72px, the logo area is 56×56 and
// centered in that rail, the same 36×36 mark sits in the middle of the
// area, and the wordmark is not shown.
//
// The rail and the logo area are the existing "sidebar" and "logo-slot"
// labels. ui.NewTester can read a labeled box; the 56×56 area has no fill
// of its own, so a pixel sample cannot see its edges.
func TestSidebarLogoSizes(t *testing.T) {
	const (
		markSize    float32 = 36
		railWidth   float32 = 72
		logoBox     float32 = 56
		wordmarkGap float32 = 5 // margin-right, spec says about 5px
	)

	m := NewModel(demoHome())
	expanded := renderModel(m, WindowWidth, WindowHeight, false)

	mark, ok := expanded.Find("mygo")
	if !ok {
		t.Fatal("expanded sidebar has no favicon mark")
	}
	if mark.W != markSize || mark.H != markSize {
		t.Fatalf("expanded mark %v, want %v×%v", mark, markSize, markSize)
	}
	name, ok := expanded.Find("Clash Go")
	if !ok {
		t.Fatal("expanded sidebar is missing the Clash Go wordmark")
	}
	gap := name.X - (mark.X + mark.W)
	if !near(gap, wordmarkGap, 1) {
		t.Fatalf("wordmark gap %v, want about %v (mark %v name %v)", gap, wordmarkGap, mark, name)
	}
	if name.Y+name.H <= mark.Y || mark.Y+mark.H <= name.Y {
		t.Fatalf("wordmark is not on the mark's row: mark %v name %v", mark, name)
	}
	if _, shown := expanded.Find("logo-slot"); shown {
		t.Fatal("expanded sidebar should not use the collapsed 56×56 logo area")
	}

	collapsed := NewCollapsedTester(m, false)

	rail, ok := collapsed.Find("sidebar")
	if !ok {
		t.Fatal("collapsed sidebar has no sidebar box")
	}
	if rail.X != 0 || rail.W != railWidth {
		t.Fatalf("collapsed rail %v, want x=0 width %v", rail, railWidth)
	}
	slot, ok := collapsed.Find("logo-slot")
	if !ok {
		t.Fatal("collapsed sidebar has no 56×56 logo area")
	}
	if slot.W != logoBox || slot.H != logoBox {
		t.Fatalf("collapsed logo area %v, want %v×%v", slot, logoBox, logoBox)
	}
	// The rail's 1px divider is inside the 72, so the gaps beside a
	// centered 56px box differ by that 1px.
	left := slot.X - rail.X
	right := (rail.X + rail.W) - (slot.X + slot.W)
	if !near(left, right, 1) || left < 0 || right < 0 {
		t.Fatalf("logo area is not centered in the rail: rail %v slot %v", rail, slot)
	}
	cmark, ok := collapsed.Find("mygo")
	if !ok {
		t.Fatal("collapsed sidebar lost the favicon mark")
	}
	if cmark.W != markSize || cmark.H != markSize {
		t.Fatalf("collapsed mark %v, want %v×%v inside the logo area", cmark, markSize, markSize)
	}
	if !containsRect(slot, cmark) {
		t.Fatalf("mark %v is outside the logo area %v", cmark, slot)
	}
	mx, my := rectCenter(cmark)
	sx, sy := rectCenter(slot)
	if !near(mx, sx, 0.5) || !near(my, sy, 0.5) {
		t.Fatalf("mark center (%v,%v), logo area center (%v,%v)", mx, my, sx, sy)
	}
	if collapsed.HasText("Clash Go") {
		t.Fatal("collapsed sidebar still shows the Clash Go wordmark")
	}
}

func rectCenter(r mygoui.Rect) (float32, float32) {
	return r.X + r.W/2, r.Y + r.H/2
}

func containsRect(outer, inner mygoui.Rect) bool {
	const eps float32 = 0.01
	return inner.X >= outer.X-eps && inner.Y >= outer.Y-eps &&
		inner.X+inner.W <= outer.X+outer.W+eps &&
		inner.Y+inner.H <= outer.Y+outer.H+eps
}

func near(got, want, tol float32) bool {
	d := got - want
	if d < 0 {
		d = -d
	}
	return d <= tol
}
