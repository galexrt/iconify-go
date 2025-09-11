package iconifygo

import "strings"

type Option = func(s *IconifyServer) error

// WithPreloadIconsets preloads the given iconsets into the server's cache.
func WithPreloadIconsets(preloadIconsets []string) Option {
	return func(s *IconifyServer) error {
		if len(preloadIconsets) > 0 {
			for _, key := range preloadIconsets {
				key = strings.TrimSuffix(key, ".json")
				iconSet, err := loadIconSet(key+".json", s.IconsetPath)
				if err != nil {
					return err
				}
				s.cache.Store(key, iconSet)
			}
		}
		return nil
	}
}

// WithHandlers sets the enabled handlers for the server.
func WithHandlers(handlers ...string) Option {
	return func(s *IconifyServer) error {
		s.Handlers = parseHandlerFlags(handlers)
		return nil
	}
}
