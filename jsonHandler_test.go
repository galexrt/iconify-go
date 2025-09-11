package iconifygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var set = IconSetResponse{
	Prefix:       "tabler",
	LastModified: 123456,
	Icons: map[string]Icon{
		"test1": {
			Body:   "<g>my body</g>",
			Width:  24,
			Height: 36,
		},
		"test2": {
			Body:  "<g>my body2</g>",
			Width: 24,
		},
		"test3": {
			Body:   "<g>my body2</g>",
			Height: 48,
		},
		"test-alias": {
			Body: "<g>my body</g>",
		},
	},
	Aliases: map[string]Alias{
		"alias1": {
			Parent: "test-alias",
			HFlip:  true,
		},
		"alias2": {
			Parent: "not-found",
		},
	},
	Width:  24,
	Height: 24,
}

func Test_parseIconSet(t *testing.T) {
	tests := []struct {
		name      string
		iconSet   IconSetResponse
		iconNames []string
		want      IconSetResponse
	}{
		{
			name:      "ValidIconSet",
			iconSet:   set,
			iconNames: []string{"test1", "test2", "test3", "not-found", "alias1", "alias2"},
			want: IconSetResponse{
				Prefix:       "tabler",
				LastModified: 123456,
				Icons: map[string]Icon{
					"test1": {
						Body:   "<g>my body</g>",
						Width:  24,
						Height: 36,
					},
					"test2": {
						Body:   "<g>my body2</g>",
						Width:  24,
						Height: 0,
					},
					"test3": {
						Body:   "<g>my body2</g>",
						Width:  0,
						Height: 48,
					},
					"test-alias": {
						Body: "<g>my body</g>",
					},
				},
				Aliases: map[string]Alias{
					"alias1": {
						Parent: "test-alias",
						HFlip:  true,
					},
				},
				NotFound: []string{"not-found", "alias2"},
				Width:    24,
				Height:   24,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseIconSet(tt.iconSet, tt.iconNames)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseIconSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseIconSet(b *testing.B) {
	iconSet, err := loadIconSet("websymbol.json", "test")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	num := 10000

	for i := 0; i < num; i++ {
		parseIconSet(iconSet, []string{"video", "find"})
	}
}
