package iconifygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestNewIconifyServer(t *testing.T) {
	tests := []struct {
		name        string
		basePath    string
		iconsetPath string
		handlers    []string
		want        *IconifyServer
	}{
		{name: "default handlers", basePath: "/icons", iconsetPath: "../json", handlers: []string{},
			want: &IconifyServer{BasePath: "/icons", IconsetPath: "../json", Handlers: handlerFlags{All: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIconifyServer(tt.basePath, tt.iconsetPath, tt.handlers...)

			if !assert.EqualValues(t, tt.want, got) {
				t.Errorf("NewIconifyServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
