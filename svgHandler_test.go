package iconifygo

import "testing"

func Test_parseNumberAndUnit(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    float64
		want2   string
		wantErr bool
	}{
		{
			name:    "Valid number with unit",
			str:     "10px",
			want:    10.0,
			want2:   "px",
			wantErr: false,
		},
		{
			name:    "Valid number without unit",
			str:     "10",
			want:    10.0,
			want2:   "",
			wantErr: false,
		},
		{
			name:    "Invalid number",
			str:     "abc",
			want:    0.0,
			want2:   "",
			wantErr: true,
		},
		{
			name:    "Empty string",
			str:     "",
			want:    0.0,
			want2:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := parseNumberAndUnit(tt.str)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getNumAndUnitFromString() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getNumAndUnitFromString() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("getNumAndUnitFromString() = %v, want %v", got, tt.want)
			}
			if tt.want2 != got2 {
				t.Errorf("getNumAndUnitFromString() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_buildSVG(t *testing.T) {
	tests := []struct {
		name          string
		icon          Icon
		iconSetWidth  int
		iconSetHeight int
		params        SVGparams
		want          string
		wantErr       bool
	}{
		{
			name: "Basic Icon",
			icon: Icon{
				Body:   "<g transform=\"translate(0, 0)\"></g>",
				Width:  100,
				Height: 150,
			},
			want: `<svg xmlns="http://www.w3.org/2000/svg" width="0.667em" height="1em" viewBox="0 0 100 150"><g transform="translate(0, 0)"></g></svg>`,
		},
		{
			name: "No Body",
			icon: Icon{
				Body:   "",
				Width:  100,
				Height: 100,
			},
			want: `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 100 100"></svg>`,
		},
		{
			name: "Replace Color",
			icon: Icon{
				Body:   "<g fill=\"currentColor\"></g>",
				Width:  0,
				Height: 0,
			},
			iconSetWidth:  50,
			iconSetHeight: 50,
			params: SVGparams{
				Color: "red",
			},
			want: `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 50 50"><g fill="red"></g></svg>`,
		},
		{
			name: "Invalid Width",
			icon: Icon{
				Body:   "<g fill=\"currentColor\"></g>",
				Width:  100,
				Height: 100,
			},
			params: SVGparams{
				Color: "red",
				Width: "abc",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := buildSVG(tt.icon, tt.iconSetWidth, tt.iconSetHeight, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("svgBuilder() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("svgBuilder() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("svgBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPrefixAndIcon(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  string
		want2 string
	}{
		{name: "valid path", path: "/icons/awesome.svg", want: "icons", want2: "awesome"},
		{name: "invalid path", path: "", want: "", want2: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := extractPrefixAndIcon(tt.path)
			if got != tt.want {
				t.Errorf("extractPrefixAndIcon() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("extractPrefixAndIcon() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
