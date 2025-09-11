package iconifygo

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

			return Icon{}, Alias{}, fmt.Errorf("parent not found for alias %s", iconName)
		}

		return Icon{}, Alias{}, fmt.Errorf("Icon %s not found", iconName)
	}
}

func (s *IconifyServer) getOrLoadIconSet(file string) (IconSetResponse, error) {
	iconSet, ok := s.cache.Load(file)
	if !ok {
		var err error
		iconSet, err = loadIconSet(file, s.IconsetPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return iconSet, &HTTPError{StatusCode: http.StatusNotFound, Message: "File not found"}
			}
			return iconSet, &HTTPError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
		}
		s.cache.Store(strings.TrimSuffix(file, ".json"), iconSet)
	}
	return iconSet, nil
}
