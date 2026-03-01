// Package menu provides menu data structures, loaders, i18n translation,
// and active-state management for Go web applications.
// It is rendering-agnostic — consuming projects choose how to render the menu.
package menu

// MenuItem represents a single menu entry.
// Items with no URL are typically section titles (non-clickable headings).
// Items with a URL are navigation links.
type MenuItem struct {
	Caption  string     `json:"caption"`            // Display text or translation key
	URL      string     `json:"url,omitempty"`      // Route path (empty for section titles)
	Icon     string     `json:"icon,omitempty"`     // CSS icon class (e.g. "bi bi-house")
	Children []MenuItem `json:"children,omitempty"` // Nested child items
	Active   bool       `json:"-"`                  // Whether this item is currently active (set by MarkActive)
}

// MenuData is the top-level structure for menu configuration.
// It contains the menu tree and optional i18n translations.
type MenuData struct {
	Menu []MenuItem                    `json:"menu"`            // Top-level menu items
	I18n map[string]map[string]string `json:"i18n,omitempty"` // locale → key → translation
}

// MenuConfig controls rendering behavior of the sidebar menu.
type MenuConfig struct {
	ActiveURL   string // Current page URL for active item highlighting
	HXTarget    string // HTMX swap target selector (default: "#content")
	HXPushURL   bool   // Whether to add hx-push-url="true" to links
	Collapsible bool   // Whether sections can be collapsed (Alpine.js)
	Locale      string // Current locale for i18n (default: "en_US")
}

// DefaultConfig returns a MenuConfig with sensible defaults.
func DefaultConfig() MenuConfig {
	return MenuConfig{
		HXTarget:    "#content",
		HXPushURL:   true,
		Collapsible: true,
		Locale:      "en_US",
	}
}

// Translate resolves a caption using the i18n map for the given locale.
// Falls back to the original caption if no translation is found.
func Translate(caption, locale string, i18n map[string]map[string]string) string {
	if i18n == nil || locale == "" || locale == "en_US" {
		return caption
	}
	if translations, ok := i18n[locale]; ok {
		if translated, ok := translations[caption]; ok && translated != "" {
			return translated
		}
	}
	return caption
}

// MarkActive sets the Active flag on menu items matching the current URL.
// It returns true if any item in the tree was marked active.
func MarkActive(items []MenuItem, currentURL string) bool {
	found := false
	for i := range items {
		items[i].Active = false
		if items[i].URL != "" && items[i].URL == currentURL {
			items[i].Active = true
			found = true
		}
		if len(items[i].Children) > 0 {
			if MarkActive(items[i].Children, currentURL) {
				found = true
			}
		}
	}
	return found
}

// HasActiveChild returns true if any direct child of the item is active.
func HasActiveChild(item MenuItem) bool {
	for _, child := range item.Children {
		if child.Active {
			return true
		}
	}
	return false
}
