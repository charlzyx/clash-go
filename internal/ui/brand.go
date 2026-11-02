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

// brandMarkSize is the sidebar title mark, next to the 16px app name.
const brandMarkSize float32 = 22
