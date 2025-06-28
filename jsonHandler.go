package iconifygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Icon struct {
	Body   string `json:"body"`
	Left   int    `json:"left,omitempty"`
	Top    int    `json:"top,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	HFlip  bool   `json:"hFlip,omitempty"`
	VFlip  bool   `json:"vFlip,omitempty"`
	Rotate int    `json:"rotate,omitempty"`
}

type Alias struct {
	Parent string `json:"parent"`
	HFlip  bool   `json:"hFlip,omitempty"`
	VFlip  bool   `json:"vFlip,omitempty"`
	Rotate int    `json:"rotate,omitempty"`
}
type IconSetResponse struct {
	Prefix       string           `json:"prefix"`
	LastModified uint             `json:"lastModified,omitempty"`
	Aliases      map[string]Alias `json:"aliases,omitempty"`
	Width        int              `json:"width,omitempty"`
	Height       int              `json:"height,omitempty"`
	Icons        map[string]Icon  `json:"icons"`
	NotFound     []string         `json:"not_found,omitempty"`
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
	var icons = make(map[string]Icon, len(iconNames))
	var notFound []string
	var aliases = make(map[string]Alias, len(iconNames))

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

func getIconFromSet(iconSet *IconSetResponse, iconName string) (Icon, Alias, error) {
	if _, ok := iconSet.Icons[iconName]; ok {
		return iconSet.Icons[iconName], Alias{}, nil
	} else {
		if _, ok := iconSet.Aliases[iconName]; ok {
			alias := iconSet.Aliases[iconName]
			if _, ok := iconSet.Icons[iconSet.Aliases[iconName].Parent]; ok {
				icon := iconSet.Icons[iconSet.Aliases[iconName].Parent]
				return icon, alias, nil
			}

			return Icon{}, Alias{}, fmt.Errorf("Parent not found for alias %s", iconName)
		}

		return Icon{}, Alias{}, fmt.Errorf("Icon %s not found", iconName)
	}
}

func (s *IconifyServer) handleJSON(w http.ResponseWriter, r *http.Request) {
	dirParts := strings.Split(r.URL.Path, "/")
	file := strings.Split(r.URL.Path, "/")[len(dirParts)-1]

	iconSet, err := loadIconSet(file, s.IconsetPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

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
