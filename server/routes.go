package main

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/server/zfs"
)

func newServer() http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/hello", zfs.HandleHello())
	mux.Handle("/", http.NotFoundHandler())

}
