package iconifygo

import (
	"testing"
)

func Test_parseHandlerFlags(t *testing.T) {
	tests := []struct {
		name     string
		handlers []string
		want     handlerFlags
	}{
		{name: "all handlers", handlers: []string{"all"}, want: handlerFlags{All: true}},
		{name: "svg handler", handlers: []string{"svg"}, want: handlerFlags{SVG: true}},
		{name: "json handler", handlers: []string{"json"}, want: handlerFlags{JSON: true}},
		{name: "no handlers", handlers: []string{}, want: handlerFlags{All: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseHandlerFlags(tt.handlers)
			if tt.want != got {
				t.Errorf("parseHandlerFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
