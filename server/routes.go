package main

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/server/firewall"
	"github.com/nickrobison/terraform-linux-provider/server/middleware"
	"github.com/nickrobison/terraform-linux-provider/server/zfs"
)

func newServer(zfsClient zfs.ZfsClient, firewallClient firewall.FirewallClient) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, zfsClient, firewallClient)
	var handler http.Handler = mux
	middleware.LoggingMiddleware(handler)
	return handler
}

func addRoutes(mux *http.ServeMux, zfsClient zfs.ZfsClient, firewallClient firewall.FirewallClient) {
	mux.Handle("/hello", zfs.HandleHello())
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/zfs/zpools", zfs.HandleZpoolList(zfsClient))
	mux.Handle("GET /firewall/zones", firewall.HandleZoneList(firewallClient))
	mux.Handle("GET /firewall/zones/{name}", firewall.HandleZoneGet(firewallClient))
	mux.Handle("POST /firewall/zones", firewall.HandleZoneCreate(firewallClient))
	mux.Handle("DELETE /firewall/zones/{name}", firewall.HandleZoneDelete(firewallClient))
	mux.Handle("POST /firewall/rules", firewall.HandleRuleAdd(firewallClient))
	mux.Handle("DELETE /firewall/rules", firewall.HandleRuleRemove(firewallClient))
}
