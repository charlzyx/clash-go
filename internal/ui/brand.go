package ui

import (
	_ "embed"

	mygoui "github.com/egoist/mygo/ui"
)

//go:embed assets/mygo-favicon.svg
var mygoFaviconSVG []byte

// brandMark is the mygo favicon: a cyan rounded square and a white ring.
// Parsed once. Image shows its own colors; Icon would flatten them to the
// text color.
var brandMark = mygoui.MustParseSVG(mygoFaviconSVG)

// Sidebar and logo slot, matching Clash Verge's layout-sidebar.
const (
	sidebarExpandedWidth  float32 = 200
	sidebarCollapsedWidth float32 = 72
	// Expanded logo row: ~58px tall, 20px horizontal padding.
	logoRowHeight   float32 = 58
	logoRowPaddingX float32 = 20
	// The mark is 36×36 with 5px of room before the wordmark.
	brandMarkSize   float32 = 36
	brandMarkMargin float32 = 5
	// Collapsed logo area is 56×56. The mark stays 36×36, centered, and the wordmark is hidden.
	logoCollapsedBox float32 = 56
)
