package iconifygo

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

type SVGparams struct {
	Color, Width, Height string
	Flip                 []string
	Rotate               string
}
