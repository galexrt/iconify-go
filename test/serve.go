//go:build exclude

package main

import (
	"net/http"

	iconifygo "github.com/andyburri/iconify-go"
)

func main() {
	mux := http.NewServeMux()
	iconify := iconifygo.NewIconifyServer("/icons", ".")
	mux.HandleFunc("/icons/", iconify.HandlerFunc())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
