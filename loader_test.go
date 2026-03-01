package menu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoader_Success(t *testing.T) {
	// Create a temp JSON file
	menuJSON := `{
		"menu": [
			{
				"caption": "Main",
				"icon": "bi bi-house",
				"children": [
					{"caption": "Customers", "url": "/customer", "icon": "bi bi-people-fill"},
					{"caption": "Actors", "url": "/actor"}
				]
			}
		],
		"i18n": {
			"pt_BR": {"Main": "Principal", "Customers": "Clientes"}
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "menu.json")
	if err := os.WriteFile(tmpFile, []byte(menuJSON), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(tmpFile)
	data, err := loader.Load()
	if err != nil {
		t.Fatalf("FileLoader.Load() error: %v", err)
	}

	if len(data.Menu) != 1 {
		t.Fatalf("expected 1 top-level item, got %d", len(data.Menu))
	}
	if data.Menu[0].Caption != "Main" {
		t.Errorf("top-level caption = %q, want %q", data.Menu[0].Caption, "Main")
	}
	if len(data.Menu[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(data.Menu[0].Children))
	}
	if data.Menu[0].Children[0].URL != "/customer" {
		t.Errorf("first child URL = %q, want %q", data.Menu[0].Children[0].URL, "/customer")
	}

	// Check i18n
	if data.I18n == nil {
		t.Fatal("i18n should not be nil")
	}
	if data.I18n["pt_BR"]["Customers"] != "Clientes" {
		t.Errorf("i18n pt_BR Customers = %q, want %q", data.I18n["pt_BR"]["Customers"], "Clientes")
	}
}

func TestFileLoader_NotFound(t *testing.T) {
	loader := NewFileLoader("/nonexistent/menu.json")
	_, err := loader.Load()
	if err == nil {
		t.Error("FileLoader.Load() should fail for nonexistent file")
	}
}

func TestFileLoader_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(tmpFile, []byte("{invalid json}"), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(tmpFile)
	_, err := loader.Load()
	if err == nil {
		t.Error("FileLoader.Load() should fail for invalid JSON")
	}
}

func TestEmbedLoader_Success(t *testing.T) {
	data := []byte(`{
		"menu": [
			{"caption": "Projects", "children": [
				{"caption": "Reviews", "url": "/review"}
			]}
		]
	}`)

	loader := NewEmbedLoader(data)
	menuData, err := loader.Load()
	if err != nil {
		t.Fatalf("EmbedLoader.Load() error: %v", err)
	}

	if len(menuData.Menu) != 1 {
		t.Fatalf("expected 1 top-level item, got %d", len(menuData.Menu))
	}
	if menuData.Menu[0].Caption != "Projects" {
		t.Errorf("top-level caption = %q, want %q", menuData.Menu[0].Caption, "Projects")
	}
}

func TestAPILoader_Success(t *testing.T) {
	menuData := MenuData{
		Menu: []MenuItem{
			{
				Caption: "Main",
				Children: []MenuItem{
					{Caption: "Customers", URL: "/customer"},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Authorization") != "Token test123" {
			t.Errorf("expected Authorization header, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept header, got %q", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(menuData)
	}))
	defer server.Close()

	loader := NewAPILoader(server.URL, map[string]string{
		"Authorization": "Token test123",
	})

	data, err := loader.Load()
	if err != nil {
		t.Fatalf("APILoader.Load() error: %v", err)
	}

	if len(data.Menu) != 1 {
		t.Fatalf("expected 1 top-level item, got %d", len(data.Menu))
	}
	if data.Menu[0].Children[0].Caption != "Customers" {
		t.Errorf("child caption = %q, want %q", data.Menu[0].Children[0].Caption, "Customers")
	}
}

func TestAPILoader_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	loader := NewAPILoader(server.URL, nil)
	_, err := loader.Load()
	if err == nil {
		t.Error("APILoader.Load() should fail for 500 status")
	}
}

func TestAPILoader_InvalidURL(t *testing.T) {
	loader := NewAPILoader("http://localhost:99999/invalid", nil)
	_, err := loader.Load()
	if err == nil {
		t.Error("APILoader.Load() should fail for unreachable URL")
	}
}
