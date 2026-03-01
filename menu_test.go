package menu

import (
	"testing"
)

func TestTranslate_EnglishDefault(t *testing.T) {
	i18n := map[string]map[string]string{
		"pt_BR": {"Customers": "Clientes"},
	}
	// en_US should return caption as-is
	got := Translate("Customers", "en_US", i18n)
	if got != "Customers" {
		t.Errorf("Translate(en_US) = %q, want %q", got, "Customers")
	}
}

func TestTranslate_PortugueseBR(t *testing.T) {
	i18n := map[string]map[string]string{
		"pt_BR": {"Customers": "Clientes", "Add": "Adicionar"},
	}
	got := Translate("Customers", "pt_BR", i18n)
	if got != "Clientes" {
		t.Errorf("Translate(pt_BR, Customers) = %q, want %q", got, "Clientes")
	}
	got = Translate("Add", "pt_BR", i18n)
	if got != "Adicionar" {
		t.Errorf("Translate(pt_BR, Add) = %q, want %q", got, "Adicionar")
	}
}

func TestTranslate_MissingKey(t *testing.T) {
	i18n := map[string]map[string]string{
		"pt_BR": {"Customers": "Clientes"},
	}
	// Missing key should fall back to caption
	got := Translate("Unknown", "pt_BR", i18n)
	if got != "Unknown" {
		t.Errorf("Translate(pt_BR, Unknown) = %q, want %q", got, "Unknown")
	}
}

func TestTranslate_MissingLocale(t *testing.T) {
	i18n := map[string]map[string]string{
		"pt_BR": {"Customers": "Clientes"},
	}
	got := Translate("Customers", "es_ES", i18n)
	if got != "Customers" {
		t.Errorf("Translate(es_ES, missing) = %q, want %q", got, "Customers")
	}
}

func TestTranslate_NilMap(t *testing.T) {
	got := Translate("Customers", "pt_BR", nil)
	if got != "Customers" {
		t.Errorf("Translate(nil) = %q, want %q", got, "Customers")
	}
}

func TestMarkActive(t *testing.T) {
	items := []MenuItem{
		{
			Caption: "Main",
			Children: []MenuItem{
				{Caption: "Customers", URL: "/customer"},
				{Caption: "Actors", URL: "/actor"},
			},
		},
		{
			Caption: "Projects",
			Children: []MenuItem{
				{Caption: "Reviews", URL: "/review"},
				{Caption: "Tasks", URL: "/task"},
			},
		},
	}

	found := MarkActive(items, "/review")
	if !found {
		t.Error("MarkActive should return true when URL matches")
	}

	// The review item should be active
	if !items[1].Children[0].Active {
		t.Error("Reviews should be active")
	}

	// Others should not be active
	if items[0].Children[0].Active {
		t.Error("Customers should not be active")
	}
	if items[1].Children[1].Active {
		t.Error("Tasks should not be active")
	}
}

func TestMarkActive_NoMatch(t *testing.T) {
	items := []MenuItem{
		{Caption: "Main", Children: []MenuItem{
			{Caption: "Customers", URL: "/customer"},
		}},
	}

	found := MarkActive(items, "/nonexistent")
	if found {
		t.Error("MarkActive should return false when no URL matches")
	}
}

func TestHasActiveChild(t *testing.T) {
	item := MenuItem{
		Caption: "Main",
		Children: []MenuItem{
			{Caption: "Customers", URL: "/customer", Active: false},
			{Caption: "Actors", URL: "/actor", Active: true},
		},
	}

	if !HasActiveChild(item) {
		t.Error("HasActiveChild should return true when a child is active")
	}

	// Deactivate all
	item.Children[1].Active = false
	if HasActiveChild(item) {
		t.Error("HasActiveChild should return false when no child is active")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.HXTarget != "#content" {
		t.Errorf("DefaultConfig().HXTarget = %q, want %q", cfg.HXTarget, "#content")
	}
	if !cfg.HXPushURL {
		t.Error("DefaultConfig().HXPushURL should be true")
	}
	if !cfg.Collapsible {
		t.Error("DefaultConfig().Collapsible should be true")
	}
	if cfg.Locale != "en_US" {
		t.Errorf("DefaultConfig().Locale = %q, want %q", cfg.Locale, "en_US")
	}
}
