package ui

import mygoui "github.com/egoist/mygo/ui"

// Settings icons follow the same monochrome outline as the sidebar
// (navicons.go). Icon paints them with the row's text color, so the
// artwork uses black and the button does not.
const settingsStroke = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">`

var (
	iconHelp     = mustNav(settingsStroke + `<circle cx="12" cy="12" r="9"/><path d="M9.4 9.4 a2.6 2.6 0 1 1 3.4 2.5 c-.8.3-1.3 1-1.3 1.8 V15"/><path d="M12 17.5 h.01"/></svg>`)
	iconInfo     = mustNav(settingsStroke + `<circle cx="12" cy="12" r="9"/><path d="M12 11 V17"/><path d="M12 8 h.01"/></svg>`)
	iconCopy     = mustNav(settingsStroke + `<rect x="8" y="8" width="11" height="11" rx="1.5"/><path d="M16 8 V5.5 H5.5 V16 H8"/></svg>`)
	iconLan      = mustNav(settingsStroke + `<rect x="3" y="3" width="6" height="6" rx="1"/><rect x="15" y="3" width="6" height="6" rx="1"/><rect x="9" y="15" width="6" height="6" rx="1"/><path d="M6 9 V12 H18 V9"/><path d="M12 12 V15"/></svg>`)
	iconChevron  = mustNav(settingsStroke + `<path d="M9 6 L15 12 L9 18"/></svg>`)
	iconTelegram = mustNav(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="#000" d="M21.5 4.2 L2.7 11.2 c-.8.3-.7 1.4.1 1.6 l4.8 1.5 1.8 5.6 c.2.7 1.1.8 1.5.2 l2.6-3.3 4.6 3.4 c.6.4 1.4.1 1.6-.6 l2.4-13.2 c.2-.9-.6-1.6-1.6-1.2 z M9.2 13.4 l8.6-5.3 c.2-.1.4.2.2.4 l-7.2 6.5-.3 3.2-1.3-4.8 z"/></svg>`)
	iconGitHub   = mustNav(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="#000" d="M12 2 C6.5 2 2 6.6 2 12.2 c0 4.5 2.9 8.3 6.9 9.7.5.1.7-.2.7-.5 v-1.7 c-2.8.6-3.4-1.2-3.4-1.2-.4-1.2-1.1-1.5-1.1-1.5-.9-.6.1-.6.1-.6 1 .1 1.5 1.1 1.5 1.1.9 1.6 2.4 1.1 3 .9.1-.7.4-1.1.6-1.4-2.2-.3-4.6-1.1-4.6-5 0-1.1.4-2 1-2.7-.1-.3-.4-1.3.1-2.7 0 0 .8-.3 2.8 1 a9.4 9.4 0 0 1 5 0 c1.9-1.3 2.7-1 2.7-1 .5 1.4.2 2.4.1 2.7.6.7 1 1.6 1 2.7 0 3.9-2.3 4.7-4.6 5 .4.3.7.9.7 1.9 v2.8 c0 .3.2.6.7.5 A10.2 10.2 0 0 0 22 12.2 C22 6.6 17.5 2 12 2 z"/></svg>`)
)

func settingsGear() *mygoui.SVG { return navIcon("settings") }
