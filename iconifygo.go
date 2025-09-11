package iconifygo

import (
	"net/http"
	"path"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

type IconifyServer struct {
	BasePath    string
	IconsetPath string
	Handlers    handlerFlags

	cache *xsync.Map[string, IconSetResponse]
}

type handlerFlags struct {
	All  bool
	SVG  bool
	JSON bool
}

func NewIconifyServer(basePath string, iconsetPath string, opts ...Option) (*IconifyServer, error) {
	flags := handlerFlags{}
	s := &IconifyServer{
		BasePath:    path.Clean("/" + basePath),
		IconsetPath: iconsetPath,
		Handlers:    flags,
		cache:       xsync.NewMap[string, IconSetResponse](),
	}

	if len(opts) == 0 {
		opts = append(opts, WithHandlers())
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func parseHandlerFlags(handlers []string) handlerFlags {
	h := handlerFlags{}
	if len(handlers) == 0 {
		h.All = true
		return h
	}

	for _, handler := range handlers {
		switch strings.ToLower(handler) {
		case "svg":
			h.SVG = true

		case "json":
			h.JSON = true

		case "all":
			h.All = true
		}
	}

	return h
}

func (s *IconifyServer) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(r.URL.Path, s.BasePath) {
			http.NotFound(w, r)
			return
		}

		trimmedPath := strings.TrimPrefix(r.URL.Path, s.BasePath)
		parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

		switch {
		case len(parts) == 1 && strings.HasSuffix(parts[0], ".json"):
			if s.Handlers.All || s.Handlers.JSON {
				s.handleJSON(w, r)
			} else {
				http.NotFound(w, r)
			}

		case len(parts) == 2 && strings.HasSuffix(parts[1], ".svg"):
			if s.Handlers.All || s.Handlers.SVG {
				s.handleSVG(w, r)
			} else {
				http.NotFound(w, r)
			}

		default:
			http.NotFound(w, r)
		}
	}
}
