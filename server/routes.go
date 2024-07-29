package main

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/server/middleware"
	"github.com/nickrobison/terraform-linux-provider/server/zfs"
)

func newServer(zfsClient zfs.ZfsClient) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, zfsClient)
	var handler http.Handler = mux
	middleware.LoggingMiddleware(handler)
	return handler
}

func addRoutes(mux *http.ServeMux, zfsClient zfs.ZfsClient) {
	mux.Handle("/hello", zfs.HandleHello())
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/zfs/zpools", zfs.HandleZpoolList(zfsClient))

}
