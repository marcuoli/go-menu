package menu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// MenuLoader defines how menu data is loaded from a source.
type MenuLoader interface {
	Load() (*MenuData, error)
}

// --- FileLoader ---

// FileLoader loads menu data from a JSON file on disk.
type FileLoader struct {
	path string
}

// NewFileLoader creates a loader that reads from a JSON file at the given path.
func NewFileLoader(path string) MenuLoader {
	return &FileLoader{path: path}
}

// Load reads and parses the JSON file.
func (l *FileLoader) Load() (*MenuData, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return nil, fmt.Errorf("go-menu: failed to read file %q: %w", l.path, err)
	}
	return parseMenuData(data)
}

// --- EmbedLoader ---

// EmbedLoader loads menu data from embedded JSON bytes (typically via go:embed).
type EmbedLoader struct {
	data []byte
}

// NewEmbedLoader creates a loader from embedded JSON bytes.
func NewEmbedLoader(data []byte) MenuLoader {
	return &EmbedLoader{data: data}
}

// Load parses the embedded JSON bytes.
func (l *EmbedLoader) Load() (*MenuData, error) {
	return parseMenuData(l.data)
}

// --- APILoader ---

// APILoader loads menu data from a REST API endpoint.
type APILoader struct {
	url     string
	headers map[string]string
	client  *http.Client
}

// NewAPILoader creates a loader that fetches menu JSON from a URL.
// Optional headers (e.g. Authorization) are sent with the request.
func NewAPILoader(url string, headers map[string]string) MenuLoader {
	return &APILoader{
		url:     url,
		headers: headers,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Load fetches menu data from the API endpoint.
func (l *APILoader) Load() (*MenuData, error) {
	req, err := http.NewRequest(http.MethodGet, l.url, nil)
	if err != nil {
		return nil, fmt.Errorf("go-menu: failed to create request for %q: %w", l.url, err)
	}

	for k, v := range l.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("go-menu: API request to %q failed: %w", l.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("go-menu: API returned status %d for %q", resp.StatusCode, l.url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("go-menu: failed to read API response: %w", err)
	}

	return parseMenuData(body)
}

// --- Helpers ---

func parseMenuData(data []byte) (*MenuData, error) {
	var menuData MenuData
	if err := json.Unmarshal(data, &menuData); err != nil {
		return nil, fmt.Errorf("go-menu: failed to parse JSON: %w", err)
	}
	return &menuData, nil
}
