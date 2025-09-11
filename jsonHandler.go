package iconifygo

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// HTTPError represents an error with an associated HTTP status code.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func loadIconSet(fileName, iconDir string) (IconSetResponse, error) {
	var iconSet IconSetResponse

	fullPath := filepath.Join(iconDir, fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return iconSet, err
	}

	if err := json.Unmarshal(data, &iconSet); err != nil {
		return iconSet, err
	}

	return iconSet, nil
}

func parseIconSet(iconSet IconSetResponse, iconNames []string) (result IconSetResponse) {
	icons := make(map[string]Icon, len(iconNames))
	var notFound []string
	aliases := make(map[string]Alias, len(iconNames))

	for _, iconName := range iconNames {
		icon, alias, err := getIconFromSet(&iconSet, iconName)
		if err != nil {
			notFound = append(notFound, iconName)
			continue
		}

		if alias.Parent != "" {
			aliases[iconName] = alias
			icons[alias.Parent] = icon
			continue
		}

		icons[iconName] = icon
	}

	result.Icons = icons
	result.NotFound = notFound
	result.Aliases = aliases
	result.Prefix = iconSet.Prefix
	result.LastModified = iconSet.LastModified

	if iconSet.Width > 0 {
		result.Width = iconSet.Width
	}
	if iconSet.Height > 0 {
		result.Height = iconSet.Height
	}

	return result
}

func (s *IconifyServer) handleJSON(w http.ResponseWriter, r *http.Request) {
	dirParts := strings.Split(r.URL.Path, "/")
	file := strings.Split(r.URL.Path, "/")[len(dirParts)-1]

	iconSet, err := s.getOrLoadIconSet(strings.TrimSuffix(file, ".json"))
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		// fallback for unexpected error types
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	iconNames := strings.Split(r.URL.Query().Get("icons"), ",")
	if len(iconNames) == 1 && iconNames[0] == "" {
		http.Error(w, "No icons specified", http.StatusBadRequest)
		return
	}

	result := parseIconSet(iconSet, iconNames)

	response, err := json.Marshal(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
