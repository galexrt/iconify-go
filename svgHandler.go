package iconifygo

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SVGparams struct {
	Color, Width, Height string
	Flip                 []string
	Rotate               string
}

func (s *IconifyServer) handleSVG(w http.ResponseWriter, r *http.Request) {
	prefix, iconName := extractPrefixAndIcon(r.URL.Path)
	if prefix == "" || iconName == "" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	iconSet, err := loadIconSet(prefix+".json", s.IconsetPath)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
			err = fmt.Errorf("icon set %s not found", prefix)
		}
		http.Error(w, err.Error(), status)
		return
	}

	icon, _, err := getIconFromSet(&iconSet, iconName)
	if err != nil {
		http.Error(w, "Icon not found", http.StatusNotFound)
		return
	}

	params := parseSVGParams(r)
	if params.hasUnsafeChars() {
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}

	svg, err := buildSVG(icon, iconSet.Width, iconSet.Height, params)
	if err != nil {
		http.Error(w, "Failed to build SVG", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, err = w.Write([]byte(svg))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func buildSVG(icon Icon, defaultW, defaultH int, params SVGparams) (string, error) {
	widthAttr, heightAttr, err := resolveSVGSize(params, icon, defaultW, defaultH)
	if err != nil {
		return "", err
	}

	transformIcon(&icon, params)

	w := icon.Width
	if w == 0 {
		w = defaultW
	}

	h := icon.Height
	if h == 0 {
		h = defaultH
	}

	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="%d %d %d %d">%s</svg>`,
		widthAttr, heightAttr, icon.Left, icon.Top, w, h, icon.Body,
	), nil
}

func resolveSVGSize(params SVGparams, icon Icon, defaultW, defaultH int) (string, string, error) {
	w := icon.Width
	if w == 0 {
		w = defaultW
	}
	h := icon.Height
	if h == 0 {
		h = defaultH
	}

	width := formatFloat(float64(w)/float64(h)) + "em"
	height := "1em"

	if params.Width != "" {
		width = params.Width
		if params.Height == "" {
			num, unit, err := parseNumberAndUnit(params.Width)
			if err != nil {
				return "", "", err
			}
			height = formatFloat(num*float64(h)/float64(w)) + unit
		}
	}

	if params.Height != "" {
		height = params.Height
		if params.Width == "" {
			num, unit, err := parseNumberAndUnit(params.Height)
			if err != nil {
				return "", "", err
			}
			width = formatFloat(num*float64(w)/float64(h)) + unit
		}
	}

	return width, height, nil
}

func parseNumberAndUnit(str string) (float64, string, error) {
	re := regexp.MustCompile(`^([\d.]+)([a-zA-Z%]*)$`)
	matches := re.FindStringSubmatch(str)
	if len(matches) < 2 {
		return 0, "", fmt.Errorf("invalid size: %s", str)
	}

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid number: %s", str)
	}

	unit := matches[2]
	return num, unit, nil
}

func formatFloat(num float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(num, 'f', 3, 64), "0"), ".")
}

func transformIcon(icon *Icon, params SVGparams) {
	if params.Color != "" {
		icon.Body = strings.ReplaceAll(icon.Body, "currentColor", params.Color)
	}
}

func extractPrefixAndIcon(path string) (prefix, icon string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "", ""
	}
	prefix = parts[len(parts)-2]
	icon = strings.TrimSuffix(parts[len(parts)-1], ".svg")
	return prefix, icon
}

func parseSVGParams(r *http.Request) SVGparams {
	return SVGparams{
		Color:  r.URL.Query().Get("color"),
		Width:  r.URL.Query().Get("width"),
		Height: r.URL.Query().Get("height"),
		Flip:   strings.Split(r.URL.Query().Get("flip"), ","),
		Rotate: r.URL.Query().Get("rotate"),
	}
}

func (p SVGparams) hasUnsafeChars() bool {
	return strings.ContainsAny(p.Color+":"+p.Width+":"+p.Height, `"`)
}
