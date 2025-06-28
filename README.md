# Iconify Server in Go

This is a simple implementation of the [iconify server API](https://iconify.design/docs/api/queries.html) as a Go module. It provides an HTTP handler to serve SVG icons or JSON data.

## Implemented Endpoints

| Endpoint                                                                 | Handler | Description                                                                      |
| ------------------------------------------------------------------------ | ------- | -------------------------------------------------------------------------------- |
| [/{prefix}.json?icons={icons}](https://iconify.design/docs/api/svg.html) | `json`  | Returns JSON data for the specified icons.                                       |
| [/{prefix}/{icon}.svg](https://iconify.design/docs/api/svg.html)         | `svg`   | Returns SVG image. Supports query parameters for `color`, `width`, and `height`. |

## Usage

**Instantiate the Server with**

```go
iconifygo.NewIconifyServer(basePath, iconsetPath, handlers...)
```

where `basePath` is the base path for serving icons and `iconsetPath` is the path to the directory containing the icon sets.
`handlers` is a slice of strings that specifies which handlers to enable. The default is `["all"]`, which enables all handlers.

Available handlers are `svg` and `json` or `all`. See [Implemented Endpoints](#implemented-endpoints) for the handled endpoints.

**Register the Handler like so:**

```go
http.HandleFunc("GET /icons", iconify.HandlerFunc(), "svg", "json")
```

The following example handles the `/icons` endpoint. It serves the JSON files from the `./iconsets` directory.

```go
import (
	"net/http"

	iconifygo "github.com/andyburri/iconify-go"
)

func main() {
	mux := http.NewServeMux()
	iconify := iconifygo.NewIconifyServer("/icons/", "./iconsets")
	mux.HandleFunc("/icons/", iconify.HandlerFunc())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
```

## Add API url to web component

To use the local API, the API URL has to be added as a new [`IconifyProvider`](https://iconify.design/docs/api/providers.html).

```html
<script>
  IconifyProviders = {
    "": {
      resources: ["/icons"],
    },
  };
</script>
```
