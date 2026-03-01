# go-menu

A reusable [templ](https://templ.guide)-based sidebar menu component for Go web applications using **HTMX + Alpine.js + Tailwind CSS**.

## Features

- Vertical sidebar with **section titles** (parent items) and **clickable child links**
- **JSON-driven** — menu structure defined in JSON, not hardcoded
- Multiple data sources: **file**, **embedded JSON**, or **REST API**
- **HTMX navigation**: `hx-get` + `hx-push-url` for SPA-style page transitions
- **Active item highlighting** (current URL matching)
- **Collapsible sections** (Alpine.js toggle)
- **Dark mode** support (Tailwind `dark:` variants)
- **i18n-ready**: accepts translated captions per locale

## Installation

```bash
go get github.com/marcuoli/go-menu@latest
```

## Quick Start

### 1. Define your menu in JSON

```json
{
  "menu": [
    {
      "caption": "Main",
      "icon": "bi bi-house",
      "children": [
        { "caption": "Customers", "url": "/customer", "icon": "bi bi-people-fill" },
        { "caption": "Positions", "url": "/position", "icon": "bi bi-briefcase" }
      ]
    },
    {
      "caption": "Projects",
      "icon": "bi bi-folder",
      "children": [
        { "caption": "Reviews", "url": "/review", "icon": "bi bi-journal-check" },
        { "caption": "Tasks",   "url": "/task",   "icon": "bi bi-list-task" }
      ]
    }
  ],
  "i18n": {
    "pt_BR": { "Main": "Principal", "Customers": "Clientes", "Projects": "Projetos" },
    "es_ES": { "Main": "Principal", "Customers": "Clientes", "Projects": "Proyectos" }
  }
}
```

### 2. Load and render in your Go handler

```go
package main

import (
    "log"
    "net/http"

    menu "github.com/marcuoli/go-menu"
)

func main() {
    loader := menu.NewFileLoader("config/menu.json")
    menuData, err := loader.Load()
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        menu.MarkActive(menuData.Menu, r.URL.Path)
        config := menu.MenuConfig{
            ActiveURL:   r.URL.Path,
            HXTarget:    "#content",
            HXPushURL:   true,
            Collapsible: true,
            Locale:      "en_US",
        }
        menu.Sidebar(menuData.Menu, config, menuData.I18n).Render(r.Context(), w)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Data Sources

### FileLoader
Reads from a JSON file on disk:
```go
loader := menu.NewFileLoader("config/menu.json")
data, err := loader.Load()
```

### EmbedLoader
Reads from embedded JSON (via `go:embed`):
```go
//go:embed config/menu.json
var menuJSON []byte

loader := menu.NewEmbedLoader(menuJSON)
data, err := loader.Load()
```

### APILoader
Fetches from a REST API endpoint:
```go
loader := menu.NewAPILoader(
    "http://localhost:8000/api/menu/",
    map[string]string{"Authorization": "Token abc123"},
)
data, err := loader.Load()
```

## Types

### MenuItem
```go
type MenuItem struct {
    Caption  string     `json:"caption"`            // Display text or translation key
    URL      string     `json:"url,omitempty"`      // Route path (empty = section title)
    Icon     string     `json:"icon,omitempty"`     // CSS icon class
    Children []MenuItem `json:"children,omitempty"` // Nested items
    Active   bool       `json:"-"`                  // Computed at render time
}
```

### MenuData
```go
type MenuData struct {
    Menu []MenuItem                    `json:"menu"`
    I18n map[string]map[string]string `json:"i18n,omitempty"`
}
```

### MenuConfig
```go
type MenuConfig struct {
    ActiveURL   string // Current page URL for highlighting
    HXTarget    string // HTMX swap target (default: "#content")
    HXPushURL   bool   // Add hx-push-url="true" to links
    Collapsible bool   // Alpine.js collapsible sections
    Locale      string // Current locale for i18n
}
```

## Helper Functions

| Function | Description |
|---|---|
| `Translate(caption, locale, i18n)` | Resolves a caption using the i18n map; falls back to original |
| `MarkActive(items, currentURL)` | Sets `Active` flag on items matching the URL |
| `HasActiveChild(item)` | Returns true if any direct child is active |
| `DefaultConfig()` | Returns a config with sensible defaults |

## i18n

Translations are stored in the JSON config under the `i18n` key. Each locale maps English captions to their translations:

```json
{
  "i18n": {
    "pt_BR": { "Customers": "Clientes", "Add": "Adicionar" },
    "es_ES": { "Customers": "Clientes", "Add": "Agregar" }
  }
}
```

English (`en_US`) is the default — captions are used as-is when no translation is found.

## Dark Mode

All CSS classes include Tailwind `dark:` variants. Toggle dark mode by adding/removing the `dark` class on the `<html>` element:

```html
<html class="dark">
```

## Requirements

- Go 1.22+
- [templ](https://templ.guide) v0.3+
- Tailwind CSS 4 (for styling)
- HTMX 2.x (for SPA navigation)
- Alpine.js (for collapsible sections)

## License

MIT
